package main

import (
	"fmt"
	"time"
)

func NewServiceC(serviceA ServiceAsync) *ServiceC {
	fmt.Println("Wait for service A ...")
	serviceA.Wait()
	fmt.Println("Service C can start now")
	return &ServiceC{}
}

type ServiceC struct {
	errsCh chan error
	stopCh chan struct{}
}

func (s *ServiceC) Start() (errsCh chan error, err error) {
	s.errsCh = make(chan error)
	s.stopCh = make(chan struct{})

	timer := time.NewTicker(time.Second * 3)

	go func() {
		for {
			select {
			case <-timer.C:
				fmt.Println("Service C running...")
			case <-s.stopCh:
				return
			}
		}
	}()

	return s.errsCh, nil
}

func (s *ServiceC) Stop() error {
	close(s.stopCh)
	close(s.errsCh)

	fmt.Println("Stopping service C")
	time.Sleep(time.Second * 5)
	return nil
}
