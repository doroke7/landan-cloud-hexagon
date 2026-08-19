package main

import (
	"fmt"
	"time"
)

func worker(id int, done chan<- struct{}) {
	defer func() {
		// 工作完成後送一個訊號
		done <- struct{}{}
	}()

	fmt.Printf("worker %d 開始\n", id)
	time.Sleep(time.Duration(id) * 300 * time.Millisecond)
	fmt.Printf("worker %d 完成\n", id)
}

/* 
	Channel 要做到 同步等待效果就是 在需要等待的地方取出 n 個 channel

*/

func main() {
	const workerCount = 4

	// buffer 設為 workerCount：
	// 即使 main 尚未開始接收，worker 完成時也不會卡住。
	done := make(chan struct{}, workerCount)

	for i := 1; i <= workerCount; i++ {
		go worker(i, done)
	}

	// 等待共收到 4 個「完成訊號」
	for i := 0; i < workerCount; i++ {
		<-done
	}

	fmt.Println("所有 worker 都完成了")
}