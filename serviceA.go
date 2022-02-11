package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrRestartServiceA = errors.New("please, restart the service")

type ServiceA struct {
	ready  chan struct{}
	errsCh chan error
	stopCh chan struct{}
}

func NewServiceA() *ServiceA {
	return &ServiceA{
		ready:  make(chan struct{}),
		errsCh: make(chan error),
		stopCh: make(chan struct{}),
	}
}

func (s *ServiceA) Start() (errs chan error, err error) {
	beforeReady := func() {
		time.Sleep(time.Second * 6)
		close(s.ready)
	}

	go s.restartOnError(s.errsCh, s.stopCh)
	go beforeReady()

	return s.errsCh, nil
}

func (s *ServiceA) Stop() error {
	close(s.stopCh)
	close(s.ready)

	fmt.Println("Stopping service A")
	time.Sleep(time.Second * 3)
	return nil
}

func (s *ServiceA) Wait() {
	<-s.ready
}

func (s *ServiceA) restartOnError(errsCh chan error, stopCh chan struct{}) {
	defer close(errsCh)

	timer := time.NewTicker(time.Second * 3)
	restartTimer := time.NewTimer(time.Second * 7)

	for {
		select {
		case <-timer.C:
			fmt.Println("Service A running...")
		case <-restartTimer.C:
			errsCh <- ErrRestartServiceA
			return
		case <-stopCh:
			return
		}
	}
}
