package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, heartbeat chan<- struct{}) {

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {

		select {

		case <-ticker.C:
			// 發送 heartbeat
			heartbeat <- struct{}{}

		case <-ctx.Done():
			fmt.Println("worker stopped")
			return
		}
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	heartbeat := make(chan struct{})

	go worker(ctx, heartbeat)

	// monitor
	for {

		select {

		case <-heartbeat:
			fmt.Println("worker alive")

		case <-time.After(3 * time.Second):
			fmt.Println("worker timeout")
			return
		}
	}
}
