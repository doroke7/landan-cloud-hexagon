package main

import (
	"time"
)

func main() {

	// 【同步代碼中】【無容量的 channel】
	// 一定要 recv 的goruntine 個數一定要等於 send的goruntine

	// 無容量的 channel，如果只有send => 會阻塞
	// 無容量的 channel，如果只有recv => 會dead lock

	ch := make(chan int)

	go func() {
		ch <- 1

	}()

	<-ch
	<-ch

	time.Sleep(time.Second * 10)
}
