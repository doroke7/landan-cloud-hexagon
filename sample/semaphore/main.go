package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// 信號量是一個允許 n 個併發的工具

func main() {
	var wg sync.WaitGroup
	ctx := context.Background()

	// 建立一個最大權重為 3 的信號量
	sem := semaphore.NewWeighted(3)

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()

			// 獲取信號量，每次佔用 1 個權重
			if err := sem.Acquire(ctx, 1); err != nil {
				return
			}

			fmt.Printf("🔥 官方信號量：任務 [%d] 執行中...\n", taskID)
			time.Sleep(1 * time.Second)

			// 釋放信號量
			sem.Release(1)
		}(i)
	}

	wg.Wait()
}