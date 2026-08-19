package main

import (
	"fmt"
	"sync"
	"time"
)

// CyclicBarrier 結構體
type CyclicBarrier struct {
	mu      sync.Mutex
	parties int           // 觸發屏障所需的目標總人數
	count   int           // 當前已經到達的人數
	barrier chan struct{} // 用來讓協程等待的廣播通道
}

// NewCyclicBarrier 初始化一個循環屏障，傳入需要等待的協程數量
func NewCyclicBarrier(parties int) *CyclicBarrier {
	if parties <= 0 {
		panic("parties 必須大於 0")
	}
	return &CyclicBarrier{
		parties: parties,
		barrier: make(chan struct{}),
	}
}

// Await 協程呼叫此方法表示「我到了，我開始等大家」
func (cb *CyclicBarrier) Await() {
	cb.mu.Lock()
	cb.count++

	// 暫存當前的 barrier 通道。
	// 因為最後一個人到的時候會替換掉 cb.barrier，我們必須保留舊通道的引用來進行 close 廣播。
	currentBarrier := cb.barrier

	if cb.count == cb.parties {
		// 【情況 A】最後一個人到了！
		fmt.Println("📢 [屏障提示] 所有人到齊！解鎖屏障，並重置以進入下一輪。")
		
		cb.count = 0                     // 1. 重置計數器
		cb.barrier = make(chan struct{}) // 2. 建立全新的通道（這就是 Cyclic 循環重複使用的關鍵）
		
		close(currentBarrier)            // 3. 敲響按鈕！關閉舊通道，瞬間廣播喚醒其餘所有人
		cb.mu.Unlock()
	} else {
		// 【情況 B】人還沒到齊
		cb.mu.Unlock()      // 先解鎖，讓其他隊友也能進來 Await()
		<-currentBarrier    // 卡在舊通道這裡，死等最後一個人來 close 它
	}
}

func main() {
	// 建立一個需要 3 個人集合的循環屏障
	barrier := NewCyclicBarrier(3)
	var wg sync.WaitGroup

	// 模擬 3 個併發的隊友
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// ====== 第一輪任務：準備出發 ======
			time.Sleep(time.Duration(id*200) * time.Millisecond) // 模擬每個人準備時間不同
			fmt.Printf("🏃 隊友 [%d] 到達【第一階段】集合點...\n", id)
			
			barrier.Await() // 進入屏障等待
			fmt.Printf("🚀 隊友 [%d] 成功通過第一階段，奔向第二階段！\n", id)

			// ====== 第二輪任務：重複使用同一個屏障 ======
			time.Sleep(time.Duration(id*100) * time.Millisecond) // 模擬第二階段耗時
			fmt.Printf("🏃 隊友 [%d] 到達【第二階段】集合點...\n", id)
			
			barrier.Await() // 再次進入同一個屏障等待
			fmt.Printf("🏁 隊友 [%d] 到達終點！\n", id)

		}(i)
	}

	wg.Wait()
	fmt.Println("🎉 所有階段皆已安全通關！")
}