package main

import (
	"errors"
	"fmt"
	"sync"
)

type Service interface {
	Start() (errs chan error, err error)
	Stop() (err error)
}

type ServiceAsync interface {
	Start() (errs chan error, err error)
	Stop() (err error)
	Wait()
}

type Orchestrator struct {
	terminated chan struct{}
	wg         sync.WaitGroup

	A *ServiceA
	B *ServiceB
	C *ServiceC
}

func (o *Orchestrator) Start() (err error) {
	o.terminated = make(chan struct{})
	o.A = NewServiceA()

	serviceAErrs, err := o.A.Start()
	if err != nil {
		return err
	}

	go o.watchErrors(serviceAErrs, "A")

	o.B = NewServiceB()
	serviceBErrs, err := o.B.Start()
	if err != nil {
		return err
	}

	go o.watchErrors(serviceBErrs, "B")

	o.C = NewServiceC(o.A)
	serviceCErrs, err := o.C.Start()
	if err != nil {
		return err
	}

	go o.watchErrors(serviceCErrs, "C")

	o.wg.Add(3)
	return nil
}

func (o *Orchestrator) Stop() (err []error) {
	defer close(o.terminated)

	errLock := make(chan struct{})
	var errs []error = make([]error, 0)

	go func() {
		defer func() {
			o.wg.Done()
			fmt.Println("SERVICE A CLOSED")
		}()

		err := o.A.Stop()
		if err != nil {
			errLock <- struct{}{}
			errs = append(errs, err)
			<-errLock
		}
	}()

	go func() {
		defer func() {
			o.wg.Done()
			fmt.Println("SERVICE B CLOSED")
		}()

		err := o.B.Stop()
		if err != nil {
			errLock <- struct{}{}
			errs = append(errs, err)
			<-errLock
		}
	}()

	go func() {
		defer func() {
			o.wg.Done()
			fmt.Println("SERVICE C CLOSED")
		}()

		err := o.C.Stop()
		if err != nil {
			errLock <- struct{}{}
			errs = append(errs, err)
			<-errLock
		}
	}()

	o.wg.Wait()
	return errs
}

func (o *Orchestrator) watchErrors(errs chan error, svc string) {
	defer func() {
		fmt.Println("WATCH ERRORS CLOSED " + svc)
	}()

	for err := range errs {
		fmt.Printf("[ERR] %s\n", err)

		if errors.Is(err, ErrCritialProblem) {
			o.Stop()
		}

		if errors.Is(err, ErrRestartServiceA) {
			fmt.Println("RESTARTING SERVICE A")
			o.A = NewServiceA()
			serviceAErrs, err := o.A.Start()

			if err != nil {
				fmt.Printf("[ERR] cannot restart service A: %v\n", err)
				o.Stop()
			}

			go o.watchErrors(serviceAErrs, "A")
		}
	}
}
