package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrRestartServiceA = errors.New("please, restart the service")

type ServiceA struct {
	ready chan struct{}
}

func NewServiceA() *ServiceA {
	return &ServiceA{
		ready: make(chan struct{}),
	}
}

func (s *ServiceA) Run(stop chan struct{}, errs chan error) {
	go s.restartOnError(stop, errs)

	time.Sleep(time.Second * 6)
	close(s.ready)

	<-stop
	fmt.Println("Stopping service A ...")
}

func (s *ServiceA) Wait() {
	<-s.ready
}

func (s *ServiceA) restartOnError(stop chan struct{}, errs chan error) {
	defer close(errs)
	timer := time.NewTicker(time.Second * 3)
	restartTimer := time.NewTimer(time.Second * 11)

	for {
		select {
		case <-timer.C:
			fmt.Println("Service A running...")
		case <-stop:
			return
		case <-restartTimer.C:
			errs <- ErrRestartServiceA
			return
		}
	}
}
