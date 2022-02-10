package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrCritialProblem = errors.New("critial problem system shutdown")

type ServiceB struct{}

func (s *ServiceB) Start(ctx context.Context) (errsCh chan error) {
	errCh := make(chan error)
	go s.critialError(ctx, errCh)

	return errCh
}

func (s *ServiceB) Stop() error { return nil }

func (s *ServiceB) critialError(ctx context.Context, errCh chan error) {
	timer := time.NewTicker(time.Second * 3)
	errTimer := time.NewTicker(time.Second * 30)

	for {
		select {
		case <-timer.C:
			fmt.Println("Service 02 running...")
		case <-errTimer.C:
			errCh <- ErrCritialProblem
		case <-ctx.Done():
			close(errCh)
			return
		}
	}
}
