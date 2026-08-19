package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	oContext, fCancel := context.WithCancel(context.Background())

	go func(oContext context.Context) {
		for {
			select {
			case <-oContext.Done():
				fmt.Println("goroutine 取消")
				return
			default:
				fmt.Println("running")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}(oContext)

	time.Sleep(10 * time.Second)
	fCancel() // 取消 goroutine
	time.Sleep(1 * time.Second)
}
