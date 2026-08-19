package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis_rate/v10" // 💡 官方正宗限流擴充包
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// 1. 初始化標準 Redis 連線
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	// 2. 直接建立官方限流器實例，免寫 Lua 腳本！
	limiter := redis_rate.NewLimiter(rdb)

	// 💡 參數定義：【5 秒內】最多放行 【10 次】
	// 用 Limit 結構體宣告，語意極度清晰！
	rateRule := redis_rate.Limit{
		Rate:   10,              // 允許次數
		Period: 5 * time.Second, // 窗口時間
	}

	limitKey := "rate:api:submit_order"
	var wg sync.WaitGroup

	fmt.Println("🎬 5 秒內瞬間湧入 15 個送出訂單的請求...")

	// 3. 模擬 15 個並發請求暴衝
	for i := 1; i <= 15; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			// ⚡ 直接調用 Allow 檢查
			res, err := limiter.Allow(ctx, limitKey, rateRule)
			if err != nil {
				fmt.Printf("⚠️ 系統錯誤: %v\n", err)
				return
			}

			// 💡 res.Allowed 是「這一刻還能放行幾個」，0 代表被限流擋下
			if res.Allowed == 0 {
				// 🚫 失敗 Action：第 11~15 個請求會直接被彈回
				// res.RetryAfter 還會貼心地告訴你還要等多久才可以再戳
				fmt.Printf("❌ 請求 %d：限流攔截！請等待 %v 後再試\n", requestID, res.RetryAfter)
				return
			}

			// 🟢 成功 Action：前 10 個請求順利通過
			fmt.Printf("🔥 請求 %d：通過限流！[Action] 正在扣減商品庫存並建立訂單...\n", requestID)

		}(i)
	}

	wg.Wait()
}
