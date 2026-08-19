package main

import (
	"fmt"
	"sync"
	"time"
)

// 模擬 Worker：每個 Worker 擁有自己獨立的輸出 Channel
func worker(id int, jobs <-chan int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for job := range jobs {
			// 模擬耗時處理
			time.Sleep(300 * time.Millisecond)
			out <- fmt.Sprintf("Worker %d 完成任務 %d", id, job)
		}
	}()
	return out
}

// 【Fan-In 核心函數】：將多個 <-chan string 合併為單一 <-chan string
func fanIn(channels ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	multiplexedStream := make(chan string)

	// 為每一個輸入 channel 啟動一個轉發 Goroutine
	output := func(c <-chan string) {
		defer wg.Done()
		for v := range c {
			multiplexedStream <- v // 轉發到統一的 channel
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	// 啟動一個 Goroutine 等待所有 input channel 結束後關閉總輸出
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func main() {
	jobs := make(chan int, 10)

	// 1. Fan-Out：啟動 3 個 Worker，得到 3 個獨立的結果通道
	ch1 := worker(1, jobs)
	ch2 := worker(2, jobs)
	ch3 := worker(3, jobs)

	// 2. 生產任務
	for i := 1; i <= 6; i++ {
		jobs <- i
	}
	close(jobs)

	// 3. Fan-In：將 3 個 Channel 合併成一個 Channel
	merged := fanIn(ch1, ch2, ch3)

	// 4. 單一消費者讀取最終合併結果
	for result := range merged {
		fmt.Println("收到結果:", result)
	}

	fmt.Println("所有結果已成功 Fan-In 並處理完畢！")
}
