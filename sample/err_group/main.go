package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	// 1. 初始化 ErrGroup。
	// 它會回傳一個 g 物件，以及一個內建「連帶取消機制」的 ctx。
	g, ctx := errgroup.WithContext(context.Background())

	fmt.Println("🚀 [主程序] 開始多路並發抓取資料...")

	// 2. 啟動第一個任務：抓取用戶資料
	g.Go(func() error {
		select {
		case <-ctx.Done(): // 執行前或執行中先檢查有沒有隊友卡死或爆炸
			return ctx.Err()
		case <-time.After(500 * time.Millisecond): // 模擬耗時 500ms
			fmt.Println("👤 [用戶服務] 資料抓取成功！")
			return nil
		}
	})

	// 3. 啟動第二個任務：抓取訂單資料
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(600 * time.Millisecond): // 模擬耗時 600ms
			fmt.Println("📦 [訂單服務] 資料抓取成功！")
			return nil
		}
	})

	// 4. 啟動第三個任務：抓取優惠券資料 (這個會故意失敗，且跑最快)
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond): // 只要 200ms 就率先爆炸
			fmt.Println("❌ [優惠券服務] 伺服器連線中斷！")
			return errors.New("Coupon Service 404 Error")
		}
	})

	// 5. 阻塞等待所有協程結束
	// g.Wait() 具備「First Error Wins」特性，它會回傳第一個發生錯誤的協程丟出來的 error
	if err := g.Wait(); err != nil {
		fmt.Printf("\n💥 [最終結果] 任務宣告失敗！擷取到的錯誤是: %v\n", err)
	} else {
		fmt.Println("\n🎉 [最終結果] 恭喜！所有資料皆順利抓取完畢！")
	}
}