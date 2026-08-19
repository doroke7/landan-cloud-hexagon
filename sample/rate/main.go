package main

import (
	"fmt"
	"sync"

	"golang.org/x/time/rate"
)

func main() {
	// 💡 這樣看就懂了：每秒只長出 5 個令牌，但桶子容量是 10 個！
	// 也就是說：平時每秒只能過 5 個請求，但第一波「突發流量」最多可以瞬間放行 10 個。
	limiter := rate.NewLimiter(5, 10)

	var wg sync.WaitGroup

	// 模擬瞬間衝進來 15 個並發請求
	fmt.Println("🎬 突發 200 個請求暴衝進來...")
	for i := 1; i <= 200; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// ⚡ 限流 Action
			if !limiter.Allow() {
				// 🚫 失敗 Action：第 11~15 個請求會因為超過桶子上限（10），直接被彈飛
				fmt.Printf("❌ 請求 %d：被限流攔截\n", id)
				return
			}

			// 🟢 成功 Action：前 10 個請求會因為桶子容量夠大（Burst=10），全部瞬間安全放行！
			fmt.Printf("🟢 請求 %d：突發放行成功！\n", id)

		}(i)
	}

	wg.Wait()
}
