package main

import (
	"fmt"
	"time"
)

func main() {

	// 簡單的說，無容量的 channel
	// 一定要 recv 的goruntine 個數一定要小於等於 send的goruntine

	// 無容量的 channel，如果只有send => 會阻塞
	// 無容量的 channel，如果只有recv => 會dead lock

	oChannel1 := make(chan int)
	oChannel2 := make(chan int)

	go func() {
		for i := 1; i <= 10; i = i + 1 {
			oChannel1 <- i
		}
	}()

	go func() {
		for i := 101; i <= 110; i = i + 1 {
			oChannel2 <- i
		}
	}()

	var iV int
	select {
	case iV = <-oChannel1:
		fmt.Println(iV)
	case iV = <-oChannel2:
		fmt.Println(iV)
	}

	time.Sleep(time.Second * 10)
}
