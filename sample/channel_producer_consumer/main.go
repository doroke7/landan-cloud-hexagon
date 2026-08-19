package main

import (
	"fmt"
	"time"
)

// 生產者：產生資料並放入 channel
func producer(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		fmt.Printf("[生產者] 正在生產資料 %d...\n", i)
		ch <- i                            // 若 channel 滿了，這裡會自動阻塞等待
		time.Sleep(500 * time.Millisecond) // 模擬生產消耗的時間
	}
	close(ch) // 生產結束，必須關閉 channel
	fmt.Println("[生產者] 所有資料已生產完畢，關閉 Channel。")
}

// 消費者：從 channel 取出資料並處理
func consumer(ch <-chan int) {
	// 使用 for range 讀取 channel，當 channel 被 close 且裡面的資料被拿光後會自動結束迴圈
	for data := range ch {
		fmt.Printf("  └─ [消費者] 成功消費資料: %d\n", data)
		time.Sleep(1 * time.Second) // 模擬消費處理較慢的情境
	}
	fmt.Println("[消費者] 無資料可消費，結束工作。")
}

func main() {
	// 建立一個容量為 2 的 Buffered Channel
	ch := make(chan int, 2)

	// 啟動生產者與消費者 Goroutine
	go producer(ch)
	go consumer(ch)

	// 簡單等待 Goroutine 執行完畢（實務上建議使用 sync.WaitGroup）
	time.Sleep(7 * time.Second)
}
