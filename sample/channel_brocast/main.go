package main

import (
	"fmt"
	"time"
)

func worker(id int, done <-chan struct{}) {
	<-done

	fmt.Println("worker", id, "收到停止訊號")
}

func main() {
	done := make(chan struct{})

	for i := 1; i <= 10; i++ {
		go worker(i, done)
	}

	time.Sleep(time.Second)

	fmt.Println("broadcast cancel")

	// close 會對所有 等待收 channel 的 程序都送出 一個消息。
	// 可以用 close 一個 len=0 的 channel 實現 廣播效果
	close(done)

	time.Sleep(time.Second)
}
