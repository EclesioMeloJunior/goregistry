package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type ServiceInterface interface {
	Start(context.Context) (errsCh chan error)
	Stop() error
}

type Registry struct {
	cancel context.CancelFunc

	services []ServiceInterface
	errChns  []chan error
}

func (r *Registry) Add(s ...ServiceInterface) {
	r.services = append(r.services, s...)
}

func (r *Registry) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	r.cancel = cancel

	for _, s := range r.services {
		serviceType := reflect.TypeOf(s)
		fmt.Printf("Starting %v\n", serviceType)
		errCh := s.Start(ctx)
		r.errChns = append(r.errChns, errCh)
	}
}

func (r *Registry) Stop() {
	for _, s := range r.services {
		serviceType := reflect.TypeOf(s)
		fmt.Printf("Stopping %v\n", serviceType)
		s.Stop()
	}
}

func (r *Registry) Notify(notifyCh chan error, whenErr error) {
	doneCh := make(chan struct{})

	for _, errCh := range r.errChns {
		go func(errCh chan error, doneCh chan struct{}) {
			for err := range errCh {
				if errors.Is(err, whenErr) {
					notifyCh <- err
					r.cancel()
				} else {
					fmt.Printf("not important err: %s\n", err.Error())
				}
			}
			doneCh <- struct{}{}
		}(errCh, doneCh)
	}

	go func() {
		<-doneCh
		close(notifyCh)
	}()
}
