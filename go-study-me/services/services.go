package services

import (
	"context"
	"fmt"
	"time"
)

// NewIdleService initializes basic service as an "idle" service -- it doesn't do anything in its Running state,
// but still supports all state transitions.
func NewIdleService(up StartingFn, down StoppingFn) *BasicService {
	run := func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	}

	return NewBasicService(up, run, down)
}

// OneIteration is one iteration of the timer service. Called repeatedly until service is stopped, or this function returns error
// in which case, service will fail.
type OneInteration func(ctx context.Context) error

// NewTimerService runs iteration function on every interval tick. When iteration returns error, service fails.
func NewTimerService(interval time.Duration, start StartingFn, iter OneInteration, stop StoppingFn) *BasicService {
	run := func(ctx context.Context) error {
		t := time.NewTicker(interval)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				err := iter(ctx)
				if err != nil {
					return err
				}
			case <-ctx.Done():
				return nil
			}
		}
	}

	return NewBasicService(start, run, stop)
}

// NewListener provides a simple way to build service listener from supplied functions.
// Functions are only called when not nil.
func NewListener(starting, running func(), stopping, terminated func(from State), failed func(from State, failure error)) Listener {
	return &funcBasedListener{
		startingFn:   starting,
		runningFn:    running,
		stoppingFn:   stopping,
		terminatedFn: terminated,
		failedFn:     failed,
	}
}

type funcBasedListener struct {
	startingFn   func()
	runningFn    func()
	stoppingFn   func(from State)
	terminatedFn func(from State)
	failedFn     func(from State, failure error)
}

func (f *funcBasedListener) Starting() {
	if f.startingFn != nil {
		f.startingFn()
	}
}

func (f *funcBasedListener) Running() {
	if f.startingFn != nil {
		f.runningFn()
	}
}

func (f *funcBasedListener) Stopping(from State) {
	if f.startingFn != nil {
		f.stoppingFn(from)
	}
}

func (f *funcBasedListener) Terminated(from State) {
	if f.startingFn != nil {
		f.terminatedFn(from)
	}
}

func (f *funcBasedListener) Failed(from State, failure error) {
	if f.startingFn != nil {
		f.failedFn(from, failure)
	}
}

func StartAndAwaitRunning(ctx context.Context, service Service) error {
	err := service.StartAsync(ctx)
	if err != nil {
		return nil
	}

	err = service.AwaitRunning(ctx)
	if e := service.FailureCase(); e != nil {
		return e
	}

	return err
}

func StopAndAwaitTerminated(ctx context.Context, service Service) error {
	service.StopAsync()
	err := service.AwaitTerminated(ctx)
	if err != nil {
		return nil
	}

	if e := service.FailureCase(); e != nil {
		return e
	}

	// can happed e.g. if context was canceled
	return err
}

func DescribeService(service Service) string {
	name := ""
	if named, ok := service.(NamedService); ok {
		name = named.ServiceName()
	}

	if name == "" {
		name = fmt.Sprintf("%v", service)
	}

	return name
}
