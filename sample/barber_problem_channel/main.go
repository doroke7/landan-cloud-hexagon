package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	// 等候室只有 3 張椅子。通道裡裝整數，代表進來排隊的顧客 ID。
	waitingChairs := make(chan int, 3)

	// 💇‍♂️ 1. 啟動理髮師執行緒（消費者）
	go func() {
		for {
			// 如果通道是空的，理髮師自動阻塞（睡覺）。
			// 有人把 ID 丟進來，理髮師就自動醒來並拿到該顧客的 ID。
			customerID := <-waitingChairs

			fmt.Printf("✂️  理髮師醒了，開始幫顧客 %d 剪髮...\n", customerID)
			time.Sleep(1 * time.Second) // 模擬剪髮時間
		}
	}()

	// 🚶 2. 模擬陸續進門的 10 位顧客（生產者）
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 隨機間隔上門
			time.Sleep(time.Duration(rand.Intn(800)) * time.Millisecond)
			fmt.Printf("🚶 顧客 %d 走進店裡...", id)

			// 核心原子操作：select 檢查通道有沒有滿
			select {
			case waitingChairs <- id:
				// 成功把自己的 ID 塞進通道 = 坐到椅子了
				fmt.Printf(" 🪑 有空位！顧客 %d 坐下等候。\n", id)
			default:
				// 通道滿了（3張椅子都有人），立刻執行 default
				fmt.Printf(" ❌ 靠！等候室滿了，顧客 %d 直接掉頭離開。\n", id)
			}
		}(i)
	}

	wg.Wait()
	time.Sleep(2 * time.Second) // 讓理髮師剪完最後剩下的人
	fmt.Println("💈 所有人剪完，下班！")
}
