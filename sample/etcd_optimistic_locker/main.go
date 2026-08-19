package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 簡單的說就是 ，利用了 etcd 自帶的 數據版本號功能，來實現樂觀鎖控制

/*
你在程式裡只需要做兩件事：
	1.	讀取時，順便把這個自帶的版本號記下來。
	2.	寫入時，把這個號碼傳回去，跟 etcd 當下的號碼比對，一樣就過，不一樣就是被插隊了。

*/

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx := context.Background()
	key := "/account/balance"

	// 初始化帳戶餘額為 500 元
	_, _ = cli.Put(ctx, key, "500")

	// 🔄 Optimistic Locking 標準流程：用 for 循環應對衝突重試
	for {
		// 1. 【Read】讀取目前餘額，並拿走當前的 ModRevision
		getResp, err := cli.Get(ctx, key)
		if err != nil {
			log.Fatalf("讀取失敗: %v", err)
		}

		kv := getResp.Kvs[0]
		currentBalance, _ := strconv.Atoi(string(kv.Value))
		mySeenRevision := kv.ModRevision // 💡 樂觀鎖的核心防偽標籤

		fmt.Printf("👀 [Read] 當前餘額: %d, 看到的版本號 (ModRevision): %d\n", currentBalance, mySeenRevision)

		// 業務邏輯判斷
		if currentBalance < 100 {
			fmt.Println("❌ 扣款失敗：餘額不足！")
			break
		}

		// 計算新餘額
		newBalance := currentBalance - 100

		// 2. 【Compare & Write】利用 etcd 的 Txn 進行原子的「比較並交換 (CAS)」
		txnResp, err := cli.Txn(ctx).
			If(clientv3.Compare(clientv3.ModRevision(key), "=", mySeenRevision)). // 檢查版本沒變
			Then(clientv3.OpPut(key, strconv.Itoa(newBalance))).                  // 沒變就寫入
			Commit()

		if err != nil {
			log.Fatalf("事務執行出錯: %v", err)
		}

		// 3. 【判斷結果】
		if txnResp.Succeeded {
			// Optimistic Locking 成功！
			fmt.Printf("🎉 [Success] 樂觀鎖校驗通過！餘額成功扣除，新餘額為: %d\n", newBalance)
			break
		} else {
			// 校驗失敗，說明有其他進程插隊，mySeenRevision 已經過期了
			fmt.Println("💥 [Conflict] 有人插隊！版本號已變更，正在重試...")
			time.Sleep(50 * time.Millisecond) // 稍微交錯並發
		}
	}
}
