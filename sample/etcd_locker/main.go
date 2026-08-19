package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// etcd-locker 本身是 etcd-mutex 語法糖包裝。本質一模一樣

// 💡 模擬一個需要高並發安全的計數器
// 在非分散式環境中，我們通常用 sync.Mutex
// 但因為我們要並行部署多台機器，所以把這裡的 Lock 介面換成 Etcd Locker
type DistributedCounter struct {
	mu sync.Locker // 👈 這裡用的是 Go 原生的標準介面！
}

func (c *DistributedCounter) Add(nodeName string) {
	fmt.Printf("🔍 %s 嘗試獲取鎖 (呼叫 Lock()，不帶 context，沒拿到就死等)...\n", nodeName)

	// ⚡ 核心 Action：加鎖
	// 雖然看起來是標準庫的 mu.Lock()，但底層其實正戳向 Etcd，排隊拿號碼牌並等待 Channel 喚醒！
	c.mu.Lock()

	defer func() {
		fmt.Printf("🔓 %s 釋放鎖 (呼叫 Unlock())...\n", nodeName)
		c.mu.Unlock() // ⚡ 解鎖 Action
	}()

	// 🟢 拿到鎖後的真實業務 Action
	fmt.Printf("👑 👑 👑 %s 成功拿到鎖！獨佔執行中...\n", nodeName)
	time.Sleep(3 * time.Second) // 模擬執行昂貴的業務 3 秒
}

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

	// 2. 建立會話
	// 💡 參數定義：租約設定為 8 秒。
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(8))
	if err != nil {
		log.Fatalf("❌ 建立會話失敗: %v", err)
	}
	defer session.Close()

	// 3. ⚡ 關鍵點：concurrency.NewLocker 直接回傳符合 sync.Locker 介面的物件
	// （*concurrency.Mutex 本身的 Lock/Unlock 是要帶 context 的，並不滿足 sync.Locker，
	//  所以不能用「先建 Mutex 再包裝」的方式，得直接呼叫 NewLocker）
	etcdLocker := concurrency.NewLocker(session, "/locks/global_counter")

	// 4. 將馬甲鎖注入到需要標準鎖的業務結構體中
	counter := &DistributedCounter{
		mu: etcdLocker, // 完美適配！
	}

	// 5. 執行任務
	counter.Add(nodeName)
}
