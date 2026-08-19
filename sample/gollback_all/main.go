package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vardius/gollback"
)

func main() {
	// 建立一個 5 秒超時的總控制 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("🚀 Gollback.Race 賽跑開始...")

	// 這裡像極了 Promise.race()，傳入一堆回傳 (interface{}, error) 的函數
	res, err := gollback.All(ctx,
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(1 * time.Second) // 慢烏龜
			return "我是任務 A（烏龜）", nil
		},
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(200 * time.Millisecond) // 快兔子
			return "我是任務 B（兔子 🏆）", nil
		},
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(500 * time.Millisecond)
			return nil, errors.New("任務 C 爆炸了")
		},
	)

	if err != nil {
		fmt.Printf("💥 賽跑出錯: %v\n", err)
	} else {
		// 預期會拿到兔子 B 的結果，因為它跑最快，而且在它跑完前，C 還沒來得及爆炸
		fmt.Printf("🏆 賽跑獲勝者結果: %v\n", res)
	}
}
