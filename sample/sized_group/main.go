package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pkgz/syncs" // 就是這個經典的三方庫！
)

// sizedGroup 是有著信號仕 思維的 waitGroup
func main() {
	// 建立一個容量限制（Size）為 3 的併發組
	// 內部一樣是幫你封裝了有緩衝的 Channel 作為 Semaphore
	swg := syncs.NewSizedGroup(3)

	fmt.Println("🚀 啟動 10 個任務，但三方庫會自動控管並發數，最多只有 3 個同時跑...")

	for i := 1; i <= 10; i++ {
		taskID := i

		// 三方庫直接把原本的 go func() 封裝成一個方法
		// 呼叫 swg.Go() 的時候，如果名額滿了（滿3人），它就會在門口阻塞排隊
		swg.Go(func(ctx context.Context) {
			fmt.Printf("📥 任務 [%d] 擠進去了，開始執行...\n", taskID)
			time.Sleep(1 * time.Second)
			fmt.Printf("✅ 任務 [%d] 搞定，釋放名額\n", taskID)
		})
	}

	// 等待所有人排隊並執行完畢
	swg.Wait()
	fmt.Println("🎉 10 個任務全部安全下班！")
}
