package main

import (
	"context"
	"fmt"
	"time"
)

type Serv03 struct{}

func (s *Serv03) Start(ctx context.Context) (errsCh chan error) {
	errsCh = make(chan error)
	timer := time.NewTicker(time.Second * 3)

	go func() {
		for {
			select {
			case <-timer.C:
				fmt.Println("Service 03 running...")
			case <-ctx.Done():
				close(errsCh)
				return
			}
		}
	}()

	return errsCh
}

func (s *Serv03) Stop() error { return nil }
