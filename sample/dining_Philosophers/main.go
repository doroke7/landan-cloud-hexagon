package main

import (
	"fmt"
	"time"
)

func philosopher(
	id int,
	left chan struct{},
	right chan struct{},
) {

	for {

		fmt.Println(
			"philosopher",
			id,
			"thinking",
		)

		time.Sleep(time.Second)

		// 拿左叉
		<-left

		fmt.Println(
			id,
			"got left fork",
		)

		// 拿右叉
		<-right

		fmt.Println(
			id,
			"eating",
		)

		time.Sleep(time.Second)

		// 放叉
		left <- struct{}{}

		right <- struct{}{}

	}
}

/*
	哲學家就餐問題實現：
	這是一個錯誤的範例。可能會併發不安全


	如何解決哲學家就餐問題：
    a. 減少哲學家數量，或限制哲學家併發 為 4個
	b. 增加 叉子數量
	c. 改變流程，如只有奇數可以運作
	d. 改變流程，給叉子加上優先級


*/

func main() {

	const n = 5

	forks := make([]chan struct{}, n)

	for i := 0; i < n; i++ {

		forks[i] = make(chan struct{}, 1)

		// 初始化叉子
		forks[i] <- struct{}{}
	}

	for i := 0; i < n; i++ {

		left := forks[i]

		right := forks[(i+1)%n]

		go philosopher(i, left, right)
	}

	time.Sleep(10 * time.Second)
}
