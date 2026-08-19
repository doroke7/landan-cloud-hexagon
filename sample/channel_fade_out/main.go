package main

import (
	"fmt"
	"sync"
	"time"
)

// 模擬 Worker（工人）
func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done() // 宣告此 Worker 完成工作 exit

	// 多個 Worker 會同時爭搶這個 jobs channel 裡面的資料
	for job := range jobs {
		fmt.Printf("[Worker %d] 取出任務 %d，開始處理...\n", id, job)
		time.Sleep(500 * time.Millisecond) // 模擬耗時運算
		fmt.Printf("  └─ [Worker %d] 任務 %d 完成！\n", id, job)
	}
}

func main() {
	const numJobs = 10
	const numWorkers = 3 // 扇出為 3 個 Worker

	jobs := make(chan int, numJobs)
	var wg sync.WaitGroup

	// 【Fan-Out 關鍵】：啟動多個 Goroutine 監聽同一個 Channel
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	// 生產者：放入 10 個任務
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // 發送完畢關閉 Channel， Worker 的 range 迴圈會自動結束

	// 等待所有 Fan-Out 的 Worker 完成
	wg.Wait()
	fmt.Println("所有任務已被平行處理完成！")
}
