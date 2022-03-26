package services

import (
	"context"
	"fmt"
	"sync"
)

type StartingFn func(serviceContext context.Context) error

type RunningFn func(serviceContext context.Context) error

type StoppingFn func(failureCase error) error

type BasicService struct {
	// functions only run, if they are not nil. If functions are nil, service will effectively do nothing
	// in given state, and go to the next one without any error
	startFn    StartingFn
	runningFn  RunningFn
	stoppingFn StoppingFn

	// everything below is protected by this mutex
	stateMu     sync.RWMutex
	state       State
	failureCase error
	listeners   []chan func(l Listener)
	serviceName string

	// closed when state reaches Running, Terminated or Failed state
	runningWaitersCh chan struct{}
	// closed when state reaches Terminated or Failed state
	terminateWaitersCh chan struct{}

	serviceContext context.Context
	serviceCancel  context.CancelFunc
}

func invalidServiceStateError(state, expected State) error {
	return fmt.Errorf("invalid service state: %v, expected: %v", state, expected)
}

func invalidServiceStateWithFailureError(state, expected State, failure error) error {
	return fmt.Errorf("invalid service state: %v, expected: %v, failure: %v", state, expected, failure)
}

func NewBasicService(start StartingFn, run RunningFn, stop StoppingFn) *BasicService {
	return &BasicService{
		startFn:            start,
		runningFn:          run,
		stoppingFn:         stop,
		state:              New,
		runningWaitersCh:   make(chan struct{}),
		terminateWaitersCh: make(chan struct{}),
	}
}

func (b *BasicService) WithName(name string) *BasicService {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()

	if b.state != New {
		return b
	}

	b.serviceName = name
	return b
}

func (b *BasicService) ServiceName() string {
	b.stateMu.RLock()
	defer b.stateMu.RUnlock()

	return b.serviceName
}

func (b *BasicService) StartAsync(parentContext context.Context) error {
	switched, oldState := b.switchState(New, Starting, func() {
		b.serviceContext, b.serviceCancel = context.WithCancel(parentContext)
		b.notifyListeners(func(l Listener) { l.Starting() }, false)
		go b.main()
	})

	if !switched {
		return invalidServiceStateError(oldState, New)
	}

	return nil
}

// Returens true, if state switch succeeds, false if it fails. Returned state is the state before switch.
// if state switching succeeds, stateFn runs with lock held.
func (b *BasicService) switchState(from, to State, stateFn func()) (bool, State) {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()

	if b.state != from {
		return false, b.state
	}

	b.state = to
	if stateFn != nil {
		stateFn()
	}

	return true, from
}

func (b *BasicService) mustSwitchState(from, to State, stateFn func()) {
	if ok, _ := b.switchState(from, to, stateFn); !ok {
		panic("switchState failed")
	}
}

func (b *BasicService) main() {
	var err error

	if b.startFn != nil {
		err = b.startFn(b.serviceContext)
	}

	if err != nil {
		b.mustSwitchState(Starting, Failed, func() {
			b.failureCase = err
			b.serviceCancel()

			close(b.runningWaitersCh)
			close(b.terminateWaitersCh)
			b.notifyListeners(func(l Listener) { l.Failed(Starting, err) }, true)
		})

		return
	}

	stoppingFrom := Starting

	if err = b.serviceContext.Err(); err != nil {
		err = nil
		goto stop
	}

	b.mustSwitchState(Starting, Running, func() {
		close(b.runningWaitersCh)
		b.notifyListeners(func(l Listener) { l.Running() }, false)
	})

	stoppingFrom = Running
	if b.runningFn != nil {
		err = b.runningFn(b.serviceContext)
	}

stop:
	failure := err
	b.mustSwitchState(stoppingFrom, Stopping, func() {
		if stoppingFrom == Starting {
			// we will not reach Running state
			close(b.runningWaitersCh)
		}
		b.notifyListeners(func(l Listener) { l.Stopping(stoppingFrom) }, false)
	})

	// Must suer we cancel the context before running stoppingFn
	b.serviceCancel()

	if b.stoppingFn != nil {
		err = b.stoppingFn(failure)
		if failure == nil {
			failure = err
		}
	}

	if failure != nil {
		b.mustSwitchState(Stopping, Failed, func() {
			b.failureCase = failure
			close(b.terminateWaitersCh)
			b.notifyListeners(func(l Listener) { l.Failed(Stopping, failure) }, true)
		})
	} else {
		b.mustSwitchState(Stopping, Terminated, func() {
			close(b.terminateWaitersCh)
			b.notifyListeners(func(l Listener) { l.Terminated(Stopping) }, true)
		})
	}
}

func (b *BasicService) StopAsync() {
	if s := b.State(); s == Stopping || s == Terminated || s == Failed {
		// no need to do anything
		return
	}

	terminated, _ := b.switchState(New, Terminated, func() {
		close(b.runningWaitersCh)
		close(b.terminateWaitersCh)
		b.notifyListeners(func(l Listener) { l.Terminated(New) }, true)
	})

	if !terminated {
		b.serviceCancel()
	}
}

func (b *BasicService) ServiceContext() context.Context {
	if s := b.State(); s == New {
		return nil
	}

	return b.serviceContext
}

func (b *BasicService) AwaitRunning(ctx context.Context) error {
	return b.awaitState(ctx, Running, b.runningWaitersCh)
}

func (b *BasicService) AwaitTerminated(ctx context.Context) error {
	return b.awaitState(ctx, Terminated, b.terminateWaitersCh)
}

func (b *BasicService) awaitState(ctx context.Context, expectedState State, ch chan struct{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		s := b.State()
		if s == expectedState {
			return nil
		}

		if failure := b.FailureCase(); failure != nil {
			return invalidServiceStateWithFailureError(s, expectedState, failure)
		}

		return invalidServiceStateError(s, expectedState)
	}
}

func (b *BasicService) FailureCase() error {
	b.stateMu.RLock()
	defer b.stateMu.RUnlock()

	return b.failureCase
}

func (b *BasicService) State() State {
	b.stateMu.RLock()
	defer b.stateMu.RUnlock()

	return b.state
}

func (b *BasicService) AddListener(listener Listener) {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()

	if b.state == Terminated || b.state == Failed {
		return
	}

	ch := make(chan func(l Listener), 4)
	b.listeners = append(b.listeners, ch)

	// each listener has its own goroutine, processing events.
	go func() {
		for lfn := range ch {
			lfn(listener)
		}
	}()
}

// lock must be held here. Read lock would be good enough, but since
// this is called from state transitions, full lock is used
func (b *BasicService) notifyListeners(lfn func(l Listener), closeChan bool) {
	for _, ch := range b.listeners {
		ch <- lfn
		if closeChan {
			close(ch)
		}
	}
}
