package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

func main() {
	ctx := context.Background()
	maxChairs := 3
	waitingCount := 0 // 目前等候室的人數

	// 2 個核心信號量
	mutex := semaphore.NewWeighted(1)                    // 互斥鎖：用來保證【檢查空位 + 坐下】的原子性
	customers := semaphore.NewWeighted(int64(maxChairs)) // 顧客訊號：容量跟等候室椅子數一致，初始為 0，用來控制理髮師睡覺與喚醒

	// 初始設定：先把 customers 消耗掉，讓理髮師一開始沒人時只能阻塞（睡覺）
	_ = customers.Acquire(ctx, int64(maxChairs))

	// 💇‍♂️ 1. 啟動理髮師（消費者）
	go func() {
		for {
			// 如果沒有顧客，理髮師就卡在這裡睡覺（計數器為 0 阻塞）
			// 有顧客進來 Release(1) 時，理髮師才會扣除 1 並醒來
			_ = customers.Acquire(ctx, 1)

			// 消費者原子操作：【叫號 + 減少等候人數】
			_ = mutex.Acquire(ctx, 1)
			waitingCount--
			fmt.Printf("✂️  理髮師叫號，一位顧客坐上理髮椅（等候室剩餘: %d）\n", waitingCount)
			mutex.Release(1)

			time.Sleep(1 * time.Second) // 模擬剪髮時間
		}
	}()

	// 🚶 2. 模擬陸續進門的 10 位顧客（生產者）
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(800)) * time.Millisecond)

			// 生產者原子操作：【檢查是否有空位 + 輸入/坐下】
			_ = mutex.Acquire(ctx, 1)
			fmt.Printf("🚶 顧客 %d 走進店裡...", id)

			if waitingCount >= maxChairs {
				// 沒空位，直接離開
				fmt.Printf(" ❌ 靠！等候室滿了，顧客 %d 直接掉頭離開。\n", id)
				mutex.Release(1)
				return
			}

			// 有空位，坐下
			waitingCount++
			fmt.Printf(" 🪑 有空位！顧客 %d 坐下等候（目前人數: %d/%d）。\n", id, waitingCount, maxChairs)

			// 坐下後，釋放一個訊號通知理髮師：「有人在等了！」
			customers.Release(1)
			mutex.Release(1)
		}(i)
	}

	wg.Wait()
	time.Sleep(2 * time.Second)
	fmt.Println("💈 所有人剪完，下班！")
}
