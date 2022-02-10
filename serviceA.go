package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type ServiceA struct{}

func (s *ServiceA) Start(ctx context.Context) (errsCh chan error) {
	errsCh = make(chan error)
	go s.restartOnError(ctx, errsCh)

	return errsCh
}

func (s *ServiceA) Stop() error { return nil }

func (s *ServiceA) restartOnError(ctx context.Context, errsCh chan error) {
	timer := time.NewTicker(time.Second * 3)
	var counter uint

	for {
		select {
		case <-timer.C:
			fmt.Printf("#%d Service 01 running...\n", counter)
			counter++
		case <-ctx.Done():
			close(errsCh)
			return
		}

		if counter == 2 {
			errsCh <- errors.New("service A got a problem")
			return
		}
	}
}
