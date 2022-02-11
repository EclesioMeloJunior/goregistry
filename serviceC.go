package main

import (
	"fmt"
	"time"
)

func NewServiceC(svcA ServiceWithWait) *ServiceC {
	return &ServiceC{
		A: svcA,
	}
}

type ServiceC struct {
	A ServiceWithWait
}

func (s *ServiceC) Run(stop chan struct{}, errs chan error) {
	fmt.Println("Wait for service A ...")
	s.A.Wait()
	fmt.Println("Service C can start now")

	timer := time.NewTicker(time.Second * 3)

	go func() {
		defer close(errs)
		for {
			select {
			case <-timer.C:
				fmt.Println("Service C running...")
			case <-stop:
				return
			}
		}
	}()

	<-stop
	fmt.Println("Stopping service C...")
}
