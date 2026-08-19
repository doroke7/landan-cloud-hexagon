package main

import (
	"fmt"
)

func main() {

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
	for iCount := 1; iCount <= 20; iCount++ {
		select {
		// 併發取值，未必是 oChannel1 先取出
		case iV = <-oChannel1:
			fmt.Println(iV)
		case iV = <-oChannel2:
			fmt.Println(iV)
		}
	}

}
