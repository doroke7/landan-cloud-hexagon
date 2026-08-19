package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// etcd-locker 本身是 etcd-mutex 語法糖包裝。本質一模一樣

func main() {
	nodeName := fmt.Sprintf("Node-PID-%d", os.Getpid())
	fmt.Printf("🎬 節點 %s 啟動...\n", nodeName)

	// 1. 初始化 Etcd 客戶端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("❌ 連線 Etcd 失敗: %v", err)
	}
	defer cli.Close()

	// 2. 建立會話 (Session)
	// 💡 參數定義：租約 TTL 設定為 6 秒
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(6))
	if err != nil {
		log.Fatalf("❌ 建立會話失敗: %v", err)
	}
	defer session.Close()

	// 3. ⚡ 核心：建立 Etcd 原生 Mutex
	// "/locks/order_payment" 是分布式鎖的全域唯一資源 Key
	mutex := concurrency.NewMutex(session, "/locks/order_payment")

	// 4. 💡 建立一個「只願意等 3 秒」的超時控制 Context
	// 這樣如果前面卡了太多節點在排隊，我們 3 秒拿不到鎖就會自動退出，不會把伺服器記憶體卡爆！
	lockCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Printf("🔍 %s 開始嘗試獲取分布式鎖 (設定 3 秒超時)...\n", nodeName)

	// 5. ⚡ 加鎖 Action：傳入帶有超時的 Context
	// 沒抢到時，Goroutine 會在底層安靜地等 Etcd 主動發 Channel 消息喚醒，直到 3 秒截止
	if err := mutex.Lock(lockCtx); err != nil {
		if err == context.DeadlineExceeded {
			// 🚫 失敗 Action A：等太久了，前面的人還沒放鎖，超時退出
			fmt.Printf("❌ %s 獲取鎖失敗：前方排隊人數過多，3 秒超時放棄！\n", nodeName)
			return
		}
		// 🚫 失敗 Action B：Etcd 網路或系統異常
		fmt.Printf("⚠️ %s 加鎖時發生系統異常: %v\n", nodeName, err)
		return
	}

	// 6. 🟢 成功 Action：順利搶到鎖，執行核心互斥業務
	fmt.Printf("👑 👑 👑 %s 成功扣留分布式鎖！開始執行扣款等關鍵 Action...\n", nodeName)

	// 模擬處理扣款業務，需要耗時 4 秒
	time.Sleep(4 * time.Second)

	// 7. ⚡ 解鎖 Action
	// 處理完必須主動放鎖，Etcd 就會「主動推消息」給下一個正在排隊等 Channel 的小弟
	if err := mutex.Unlock(context.Background()); err != nil {
		fmt.Printf("⚠️ %s 解鎖失敗: %v\n", nodeName, err)
		return
	}
	fmt.Printf("🔓 %s 業務執行完畢，成功釋放分布式鎖！\n", nodeName)
}
