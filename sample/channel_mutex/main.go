
package main

import (
	"fmt"
	"sync"
	"time"
)

/*
在 Go 語言中，要用 Channel 實現互斥鎖（Mutex）其實非常直覺！
我們只需要利用 「容量為 1 的有緩衝 Channel（Buffered Channel）」 就能輕鬆達成：
•	上鎖（Lock）：往 Channel 塞入一個元素（如果裡面已經有東西，其他協程就會被阻塞，達到等待鎖的效果）。
•	解鎖（Unlock）：從 Channel 拿出一個元素（釋放空間，讓其他阻塞的協程可以進來）。

*/

// 1. 定義我們的 Channel 鎖結構
type ChannelMutex struct {
	sem chan struct{}
}

// 初始化一個容量為 1 的 Channel 鎖
func NewChannelMutex() *ChannelMutex {
	return &ChannelMutex{
		sem: make(chan struct{}, 1), // 容量為 1 是關鍵！
	}
}

// 上鎖：誰成功塞入 struct{}{}，誰就搶到鎖
func (m *ChannelMutex) Lock() {
	m.sem <- struct{}{}
}

// 解鎖：把裡面的東西拿出來，空出位置讓別人塞
func (m *ChannelMutex) Unlock() {
	<-m.sem
}

func main() {
	var wg sync.WaitGroup
	counter := 0
	mutex := NewChannelMutex() // 建立 Channel 鎖

	// 啟動 1000 個 Goroutine，每個去加 1
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			mutex.Lock()         // 【上鎖】
			counter++            // 臨界區：安全地操作共享變數
			time.Sleep(time.Microsecond) // 稍微模擬一下微小的運算延遲
			mutex.Unlock()       // 【解鎖】
		}()
	}

	wg.Wait()
	fmt.Printf("最終計數器的值 (預期為 1000): %d\n", counter)
}