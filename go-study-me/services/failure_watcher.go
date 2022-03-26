package services

import "github.com/pkg/errors"

type FailureWatcher struct {
	ch chan error
}

func NewFailureWatcher() *FailureWatcher {
	return &FailureWatcher{ch: make(chan error)}
}

func (w *FailureWatcher) Chan() <-chan error {
	if w == nil {
		return nil
	}

	return w.ch
}

func (w *FailureWatcher) WatchService(service Service) {
	service.AddListener(NewListener(nil, nil, nil, nil, func(from State, failure error) {
		w.ch <- errors.Wrapf(failure, "service %s failed", DescribeService(service))
	}))
}
