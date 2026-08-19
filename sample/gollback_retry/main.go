package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vardius/gollback"
)

func main() {
	ctx := context.Background()

	// 模擬一個計數器，記錄這是第幾次嘗試
	attempts := 0

	fmt.Println("🚀 開始執行具備「自動重試機制」的網路請求...")

	// 我 := "1"
	// 呼叫 gollback.Retry
	// 參數 1：控制超時的 context
	// 參數 2：最大重試次數（這裡設為 5 次）
	// 參數 3：你要執行的閉包任務
	res, err := gollback.Retry(ctx, 5, func(ctx context.Context) (interface{}, error) {
		attempts++
		fmt.Printf("⏳ [第 %d 次嘗試] 正在連線到伺服器...\n", attempts)


		// 模擬前 2 次都失敗，第 3 次才成功
		if attempts < 3 {
			time.Sleep(100 * time.Millisecond) // 模擬網路卡頓
			return nil, errors.New("503 Service Unavailable (網路逾時)")
		}

		// 第 3 次成功了！
		time.Sleep(100 * time.Millisecond)
		return "🎉 連線成功！順利取得用戶資料！", nil
	})

	// 最終結果判定
	if err != nil {
		// 如果重試了 5 次通通都失敗，才會走到這裡
		fmt.Printf("\n💥 最終結果：失敗！已達最大重試次數。錯誤原因: %v\n", err)
	} else {
		// 只要在 5 次內有任何一次成功，就會立刻帶著結果跳出重試迴圈
		fmt.Printf("\n🏆 最終結果：成功！回傳資料 -> %v\n", res)
		fmt.Printf("📊 總共嘗試了 %d 次。\n", attempts)
	}
}
