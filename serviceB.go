package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrCritialProblem = errors.New("critial problem system shutdown")

func NewServiceB() *ServiceB {
	return &ServiceB{}
}

type ServiceB struct {
	errsCh chan error
	stopCh chan struct{}
}

func (s *ServiceB) Start() (errsCh chan error, err error) {
	s.errsCh = make(chan error)
	s.stopCh = make(chan struct{})

	go s.critialError()

	return s.errsCh, nil
}

func (s *ServiceB) Stop() error {
	close(s.stopCh)

	fmt.Println("Stopping service B")
	time.Sleep(time.Second * 2)
	return nil
}

func (s *ServiceB) critialError() {
	defer close(s.errsCh)

	timer := time.NewTicker(time.Second * 3)
	errTimer := time.NewTicker(time.Second * 30)

	for {
		select {
		case <-timer.C:
			fmt.Println("Service B running...")
		case <-errTimer.C:
			s.errsCh <- ErrCritialProblem
			return
		case <-s.stopCh:
			return
		}
	}
}
