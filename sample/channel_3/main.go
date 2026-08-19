package main

import (
	"fmt"
	"time"
)

func main() {

	/*
	   無緩衝(len=0) channel,

	   必須要兩個 gorutine 同時監控 send 跟received ，否則會死鎖

	*/
	c := make(chan int)

	go func() {
		c <- 100
	}()

	go func() {
		x := <-c
		fmt.Println(x)
	}()

	time.Sleep(1 * time.Second)

}
