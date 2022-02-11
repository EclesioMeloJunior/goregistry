package main

import (
	"fmt"
	"sync"
)

type Service interface {
	Run(stop chan struct{}, errs chan error)
}

type ServiceWithWait interface {
	Service
	Wait()
}

type Orchestrator struct {
	errsDone   chan struct{}
	terminated chan struct{}
	errs       []chan error

	A     *ServiceA
	stopA chan struct{}

	B     *ServiceB
	stopB chan struct{}

	C *ServiceC
}

func (o *Orchestrator) Start() (err error) {
	o.terminated = make(chan struct{})
	o.A = NewServiceA()
	o.B = NewServiceB()
	o.C = NewServiceC(o.A)

	o.stopA = make(chan struct{})
	svcAErrs := make(chan error)
	go o.A.Run(o.stopA, svcAErrs)

	o.errs = append(o.errs, svcAErrs)

	o.stopB = make(chan struct{})
	svcBErrs := make(chan error)
	go o.B.Run(o.stopB, svcBErrs)

	o.errs = append(o.errs, svcBErrs)

	svcCErrs := make(chan error)
	go o.C.Run(o.stopA, svcCErrs)

	o.errs = append(o.errs, svcCErrs)

	o.errsDone = make(chan struct{})

	errCh := make(chan error)
	go o.watchErrors(errCh)
	go o.errHandler(errCh)
	return nil
}

func (o *Orchestrator) Stop() {
	defer close(o.terminated)

	close(o.stopA)
	close(o.stopB)

	<-o.errsDone
}

func (o *Orchestrator) watchErrors(errCh chan<- error) {
	var wg sync.WaitGroup

	for _, ch := range o.errs {
		wg.Add(1)

		go func(wg *sync.WaitGroup, ch chan error) {
			defer func() {
				wg.Done()
				fmt.Println("WATCH ERROR CLOSED")
			}()

			for err := range ch {
				errCh <- err
			}
		}(&wg, ch)
	}

	go func() {
		wg.Wait()
		close(o.errsDone)
		close(errCh)
	}()
}

func (o *Orchestrator) errHandler(errs <-chan error) {
	for err := range errs {
		fmt.Println(err)

		switch err {
		case ErrCritialProblem:
			o.Stop()
		case ErrRestartServiceA:
			fmt.Println("implement restarting ...")
		}
	}
}
