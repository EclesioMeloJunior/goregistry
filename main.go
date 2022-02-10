package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	reg := Registry{
		services: []ServiceInterface{},
		errChns:  []chan error{},
	}

	reg.Add(
		&ServiceA{},
		&ServiceB{},
		&Serv03{},
	)

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	reg.Start(timeoutCtx)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)

		serviceErrCh := make(chan error)
		reg.Notify(serviceErrCh, ErrCritialProblem)

		select {
		case <-sigc:
		case err := <-serviceErrCh:
			fmt.Printf("SERVICE REGISTRY NOTIFY GOT AN ERROR: %s\n", err.Error())
		}
		reg.Stop()
	}()

	wg.Wait()
}
