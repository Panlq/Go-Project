package services

import (
	"context"
	"fmt"
)

type State int

const (
	New State = iota
	Starting
	Running
	Stopping
	Terminated
	Failed
)

func (s State) String() string {
	switch s {
	case New:
		return "New"
	case Starting:
		return "Starting"
	case Running:
		return "Running"
	case Stopping:
		return "Stopping"
	case Terminated:
		return "Terminated"
	case Failed:
		return "Failed"
	default:
		return fmt.Sprintf("Unknown state: %d", s)
	}
}

type Service interface {
	StartAsync(ctx context.Context) error

	AwaitRunning(ctx context.Context) error

	StopAsync()

	AwaitTerminated(ctx context.Context) error

	FailureCase() error

	State() State

	AddListener(listener Listener)
}

type NamedService interface {
	Service

	ServiceName() string
}

// Listener receives notifications about Service state changes.
type Listener interface {
	Starting()

	Running()

	Stopping(from State)

	Terminated(from State)

	Failed(from State, failure error)
}
