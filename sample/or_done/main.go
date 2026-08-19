package main

import (
	"fmt"
	"time"
)

// OrDone 接收一個 done 訊號與一個資料來源 ch。
// 它會回傳一個新唯讀 Channel，只要 done 被關閉，或者 ch 讀完，回傳的 Channel 就會自動關閉。
func OrDone[T any](done <-chan struct{}, ch <-chan T) <-chan T {
	valChan := make(chan T)

	go func() {
		defer close(valChan) // 確保協程結束時，關閉輸出的 Channel
		for {
			select {
			case <-done:
				// 1. 收到外部取消訊號，立刻退出
				return
			case v, ok := <-ch:
				if !ok {
					// 2. 來源 Channel 已被關閉，正常結束
					return
				}
				// 3. 將讀到的資料轉發出去
				//    這裡必須再次搭配 select，防止在轉發阻塞時，外部突然發出 done 訊號而卡死
				select {
				case valChan <- v:
				case <-done:
					return
				}
			}
		}
	}()

	return valChan
}

func main() {
	// 建立控制訊號與資料傳輸的 Channel
	done := make(chan struct{})
	dataCh := make(chan int)

	// ==========================================
	// 協程 1：模擬資料生產者 (持續產生 10 筆資料)
	// ==========================================
	go func() {
		defer close(dataCh)
		for i := 1; i <= 10; i++ {
			time.Sleep(100 * time.Millisecond) // 模擬生產耗時
			dataCh <- i
		}
	}()

	// ==========================================
	// 協程 2：模擬外部中斷 (在 350ms 後突然取消)
	// ==========================================
	go func() {
		time.Sleep(350 * time.Millisecond)
		fmt.Println("🛑 [系統通知] 外部觸發取消訊號 (Close Done)！")
		close(done) // 關閉 done 代表發送取消信號
	}()

	// ==========================================
	// 主協程：使用 Or-Done 模式優雅、安全地讀取資料
	// ==========================================
	fmt.Println("🚀 開始讀取資料流...")

	// 透過 OrDone 包裝後，我們可以直接用最乾淨的 for-range 迴圈
	// 不用在迴圈內部寫複雜的 select-case 來監聽取消訊號
	for val := range OrDone(done, dataCh) {
		fmt.Printf("👉 成功處理資料: %d\n", val)
	}

	fmt.Println("✨ 程式優雅結束，無任何 Goroutine 洩漏。")
}