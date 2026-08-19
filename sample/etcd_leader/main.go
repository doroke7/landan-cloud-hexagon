package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {
	// 1. 用當前進程的 PID 當作節點名稱，方便測試時看清是哪台機器當選
	nodeName := fmt.Sprintf("Node-PID-%d", os.Getpid())
	fmt.Printf("🎬 節點 %s 啟動，準備加入選戰...\n", nodeName)

	// 2. 初始化 Etcd 客戶端 (連線到本地或遠端的 Etcd 叢集)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("❌ 連線 Etcd 失敗: %v", err)
	}
	defer cli.Close()

	// 3. 建立帶有租約 (TTL) 的會話 (Session)
	// 💡 參數定義：租約為 5 秒。
	// Go SDK 會自動在背景幫你「續租 (KeepAlive)」。如果這台機器斷電掛了，5 秒後租約到期，
	// Etcd 就會主動通知下一個排隊的人。
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(5))
	if err != nil {
		log.Fatalf("❌ 建立會話失敗: %v", err)
	}
	defer session.Close()

	// 4. 宣告一個全域唯一的選舉「選區 (Election Key)」
	// 所有想要競爭「同一個定時任務/同一個服務」的實例，都要填寫同一個 Key
	election := concurrency.NewElection(session, "/election/my_cron_job")

	// 處理系統訊號 (Ctrl+C)，優雅退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Printf("\n🛑 %s 收到退出訊號，正在主動退選並關閉...\n", nodeName)
		cancel()
	}()

	// 5. ⚡ 核心 Action：開始競選
	fmt.Printf("🔍 %s 正在全力競選中...\n", nodeName)
	fmt.Println("💡 [底層機制]：如果前面有人，我的 Goroutine 現在會「安靜地卡在 Channel 上」，等 Etcd 主動叫號...")

	// 💡 Campaign 是一個阻塞函數：
	// - 如果現在沒人，它會立刻搶佔成功並往下走。
	// - 如果別人已經是老大，它就會在【 這裡卡住排隊】，【不耗費任何 CPU】，直到前任老大掛掉、Etcd 主動推播通知為止。
	if err := election.Campaign(ctx, nodeName); err != nil {
		fmt.Printf("❌ 競選過程中斷: %v\n", err)
		return
	}

	// 6. 👑 當選 Leader 的真實業務 Action
	fmt.Printf("\n👑 👑 👑 恭喜！%s 正式當選為 Leader！ 👑 👑 👑\n", nodeName)
	fmt.Println("📢 [Leader Action] 只有我是老大，開始獨佔跑髒活（例如分散式 Crontab）...")

	// 模擬當上 Leader 後，持續跑業務
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 只有當選的老大能看到這個輸出
			fmt.Printf("💾 [Leader Action] %s 正在安全地寫入關鍵資料庫...\n", nodeName)
		case <-ctx.Done():
			return
		}
	}
}
