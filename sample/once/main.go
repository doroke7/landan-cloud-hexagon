package main

import (
	"fmt"
	"sync"
)

func main() {
	var oSyncOne sync.Once
	var wg sync.WaitGroup
	wg.Add(10)

	// 模擬 10 個併發請求
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done()

			// 關鍵點：這段代碼雖然被呼叫 10 次
			// 但只有第一個進去的會執行內容，其他人會排隊等它執行完後直接跳過
			oSyncOne.Do(func() {
				fmt.Printf("--- [ID:%d] once 正在執行唯一的初始化邏輯 ---\n", id)
			})

			fmt.Printf("Goroutine %d 執行完畢\n", id)
		}(i)
	}

	wg.Wait()
}
