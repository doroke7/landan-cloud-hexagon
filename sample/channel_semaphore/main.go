package main

import (
	"fmt"
	"sync"
	"time"
)

// 這是一個限制「最多同時只有 3 個任務執行」的信號量範例：

func main() {
	var wg sync.WaitGroup

	// 1. 建立一個容量為 4 （代表允許的最大並發數量）的 Channel，作為我們的信號量
	oSemaphore := make(chan struct{}, 4)

	for iIndex := 1; iIndex <= 800; iIndex++ {
		wg.Add(1)

		go func(taskID int) {
			defer wg.Done()

			// 2. 獲取信號量（Acquire）：塞入一個 struct。如果滿了，就會在這裡卡住（阻塞）
			oSemaphore <- struct{}{}

			// --- 臨界區開始：保證同一時間最多只有 3 個 Goroutine 在這 ---
			fmt.Printf("🚀 任務 [%d] 開始執行 (當前並發數已 +1)\n", taskID)
			time.Sleep(1 * time.Second) // 模擬耗時的任務
			fmt.Printf("✅ 任務 [%d] 順利完成\n", taskID)
			// --- 臨界區結束 ---

			// 3. 釋放信號量（Release）：拿出一個 struct，空出位置讓別人進來
			<-oSemaphore
		}(iIndex)
	}

	// WaitGroup 是 800 個任務的終止條件
	// oSemaphore 是 4 個任務的並行條件
	wg.Wait()
	fmt.Println("🎉 所有任務皆已處理完畢！")
}
