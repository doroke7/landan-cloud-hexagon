package main

import (
	"fmt"
	"time"
)

func worker(id int, token chan struct{}, next chan struct{}) {

	for {

		// 等待 token
		<-token

		// 取得控制權
		fmt.Println(
			"worker",
			id,
			"working",
		)

		time.Sleep(time.Second)

		// 把 token 傳給下一個
		next <- struct{}{}
	}
}

func main() {

	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	ch3 := make(chan struct{})

	go worker(1, ch1, ch2)

	go worker(2, ch2, ch3)

	go worker(3, ch3, ch1)

	// 發出第一個 token
	ch1 <- struct{}{}

	time.Sleep(10 * time.Second)
}
