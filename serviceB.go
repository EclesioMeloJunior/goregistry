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
}

func (s *ServiceB) Run(stop chan struct{}, errs chan error) {
	go s.critialError(stop, errs)

	<-stop
	fmt.Println("Stopping service B...")
}

func (s *ServiceB) critialError(stop chan struct{}, errs chan error) {
	defer close(errs)

	timer := time.NewTicker(time.Second * 3)
	errTimer := time.NewTicker(time.Second * 30)

	for {
		select {
		case <-timer.C:
			fmt.Println("Service B running...")
		case <-errTimer.C:
			errs <- ErrCritialProblem
			return
		case <-stop:
			return
		}
	}
}
