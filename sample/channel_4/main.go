package main

import (
	"fmt"
)

func main() {

	/*
	   無緩衝(len=0) channel,

	   必須要兩個 gorutine 同時監控 send 跟received ，否則會死鎖

	*/

	c := make(chan int)
	c <- 100

	x := <-c
	fmt.Println(x)

}
