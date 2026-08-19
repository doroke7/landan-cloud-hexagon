package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

// 一次啟動 三個ps “： go run main.go -id=1 & go run main.go -id=2 & go run main.go -id=3

func main() {
	// 1. 透過 flag 讓每個進程在啟動時指定自己的 ID（模擬完全獨立的進程）
	processID := flag.Int("id", 0, "獨立進程的識別碼 (ID)")
	requiredProcesses := flag.Int("count", 3, "必須到齊的進程數量")
	flag.Parse()

	if *processID == 0 {
		fmt.Println("❌ 錯誤：請使用 -id 參數指定進程編號，例如：./main -id=1")
		os.Exit(1)
	}

	fmt.Printf("🎬 [進程 %d] 啟動成功，PID: %d\n", *processID, os.Getpid())

	// 2. 連線到遠端或在地的 Etcd 叢集
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, // 生產環境可改為 Etcd 叢集 IPs
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("❌ 連線 Etcd 失敗: %v", err)
	}
	defer cli.Close()

	// 3. 建立獨立的 Session（每個進程在 Etcd 裡都有自己專屬的 Lease 心跳租約）
	session, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatalf("❌ 建立 Session 失敗: %v", err)
	}
	defer session.Close()

	// 收到中斷訊號時主動關閉 session，讓 lease 立刻撤銷，
	// 避免其他還在等待的進程需要卡到 lease TTL 到期才能繼續。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		session.Close()
	}()

	// 4. 定義跨進程共享的柵欄 Key
	barrierKey := "/my-distributed-double-barrier"

	// 建立 DoubleBarrier 實例（來自 experimental/recipes，並非 concurrency package）
	db := recipe.NewDoubleBarrier(session, barrierKey, *requiredProcesses)

	fmt.Printf("🚶 [進程 %d] 到達入口柵欄，目前進程數不足 %d，開始跨網路卡住等待...\n", *processID, *requiredProcesses)

	// =================================================================
	// 🚪 第一重屏障：Enter() —— 湊齊 3 個獨立進程才准「同時」進門開工
	// =================================================================
	// 這裡完全沒有共享記憶體，底層是用各自的 Session ID 去 Etcd 註冊節點。
	if err := db.Enter(); err != nil {
		log.Fatalf("❌ 進入柵欄失敗: %v", err)
	}

	// =================================================================
	// 🛠️ 核心平行業務期（各進程在自己的 CPU 與記憶體空間獨立狂奔）
	// =================================================================
	fmt.Printf("🚀 [進程 %d] 成功突破入口！開始執行本機的耗時重度運算...\n", *processID)

	// 模擬該進程在實際生產環境中的運算時間（例如大數據處理或影片轉檔）
	time.Sleep(5 * time.Second)

	fmt.Printf("🏁 [進程 %d] 重度運算完成！到達出口柵欄，等待其他進程也做完...\n", *processID)

	// =================================================================
	// 🚪 第二重屏障：Leave() —— 湊齊 3 個獨立進程都做完，才准「同時」解散
	// =================================================================
	// 跑得快的進程不能先死掉或中斷，必須在 Leave() 卡住，確保整個分散式任務的階段一致性。
	if err := db.Leave(); err != nil {
		log.Fatalf("❌ 離開柵欄失敗: %v", err)
	}

	fmt.Printf("👋 [進程 %d] 大家都順利做完了，跨進程同步完美結束，安全下班！\n", *processID)
}
