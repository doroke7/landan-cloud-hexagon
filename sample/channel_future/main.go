package main

import (
	"fmt"
	"time"
)

type Future struct {
	result chan int
}

func NewFuture() *Future {
	return &Future{
		result: make(chan int, 1),
	}
}

// Promise fulfill
func (f *Future) Set(value int) {
	f.result <- value
	close(f.result)
}

// Future get
func (f *Future) Await() int {
	return <-f.result
}

func asyncCalculate() *Future {

	future := NewFuture()

	go func() {

		fmt.Println("開始計算...")

		time.Sleep(8 * time.Second)

		result := 100

		// 完成 Promise
		future.Set(result)

	}()

	return future
}

func main() {

	// 立即返回 Future
	future := asyncCalculate()

	fmt.Println("主程式繼續做其他事情")

	// 等待結果
	result := future.Await()

	fmt.Println("結果:", result)
}
