package main

import (
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes" // 💡 Etcd 佇列在這個擴充包裡
)

func main() {
	// 1. 初始化 Etcd 客戶端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("❌ 連線 Etcd 失敗: %v", err)
	}
	defer cli.Close()

	// 2. 宣告一個全域唯一的佇列名稱 (Queue Key)
	// 所有分散式節點只要認準這個 Key，就能一起排隊塞任務、搶任務
	q := recipe.NewQueue(cli, "/my_global_queue")

	// 🛠️ 【生產者 Goroutine】負責塞任務進隊列
	go func() {
		time.Sleep(1 * time.Second) // 故意等一下再塞
		fmt.Println("🏭 [生產者] 正在準備把任務塞進 Etcd 隊列...")

		// ⚡ 核心 Action：Enqueue (入隊列)
		_ = q.Enqueue("Task-999")
		fmt.Println("✅ [生產者] 成功將 Task-999 送入隊列！")

		_ = q.Enqueue("Task-888")
		fmt.Println("✅ [生產者] 成功將 Task-888 送入隊列！")
	}()

	// 🍽️ 【消費者】負責從隊列拿任務出來吃
	fmt.Println("🔍 [消費者] 嘗試從隊列中提取任務...")
	fmt.Println("💡 [底層機制]：如果此時隊列是空的，Consumer 會安靜地卡在 Watch Channel 上，等 Etcd 通知...")

	// ⚡ 核心 Action：Dequeue (出隊列)
	// 這是阻塞函數！隊列沒東西時它會死死卡住（不耗 CPU），直到生產者塞東西進來，它會瞬間被喚醒！
	task1, _ := q.Dequeue()
	fmt.Printf("👑 [消費者] 成功搶到並處理任務: %s\n", task1)

	task2, _ := q.Dequeue()
	fmt.Printf("👑 [消費者] 成功搶到並處理任務: %s\n", task2)
}
