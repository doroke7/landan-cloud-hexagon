package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Println(i) // 坑：這裡印出來可能全是 5，因為 Goroutine 啟動時 i 已經跑完了
		}()
	}

	time.Sleep(time.Second * 10) // 停頓，讓 goroutine 繼續跑
}
