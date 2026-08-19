package main

import (
	"fmt"
	"time"
)

func worker(done <-chan struct{}) {

	for {
		select {
		case <-done: // 模仿 context 的取消
			fmt.Println("worker cancelled")
			return

		default:
			fmt.Println("working")
			time.Sleep(time.Second)
		}
	}
}

func main() {

	done := make(chan struct{})

	go worker(done)

	time.Sleep(3 * time.Second)

	// 模仿 context 的取消 發送取消訊號
	close(done)

	time.Sleep(time.Second)
}
