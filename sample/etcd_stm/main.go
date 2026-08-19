package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

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
	accA := "/bank/account_A"
	accB := "/bank/account_B"

	// 初始化餘額
	_, _ = cli.Put(ctx, accA, "300")
	_, _ = cli.Put(ctx, accB, "100")
	fmt.Println("🚀 [初始化] A = 300, B = 100")

	// 💡 呼叫高階 STM：
	// 你完全不需要自己寫 for 迴圈，不需要記住 Revision，也不用寫一堆 clientv3.Compare
	_, err = concurrency.NewSTM(cli, func(stm concurrency.STM) error {
		fmt.Println("🔄 [STM 嘗試執行/重試中...]")

		// 1. 【Read】直接 Get。STM 會在內存裡自動幫你盯緊這兩個 Key 的版本號
		valA := stm.Get(accA)
		valB := stm.Get(accB)

		balanceA, _ := strconv.Atoi(valA)
		balanceB, _ := strconv.Atoi(valB)

		// 2. 隨便你寫複雜的 Go 業務邏輯判斷
		if balanceA < 100 {
			return fmt.Errorf("A 餘額不足，阻止交易") // 回傳 error，STM 就會直接取消，不對 etcd 產生副作用
		}

		// 3. 【Write】計算完直接 Put。這時資料只在內存緩衝區（Write Set）
		stm.Put(accA, strconv.Itoa(balanceA-100))
		stm.Put(accB, strconv.Itoa(balanceB+100))

		return nil
		// 🚪 函式結束，STM 會自動把剛才盯緊的版本號拼裝成一個底層 Txn 送給 etcd。
		// 如果這期間 A 或 B 被別人偷動過導致失敗，STM 底層會攔截，並自動重新進來這個函式再跑一次！
	})

	if err != nil {
		log.Fatalf("❌ 交易最終失敗: %v", err)
	}

	// 驗證結果
	resA, _ := cli.Get(ctx, accA)
	resB, _ := cli.Get(ctx, accB)
	fmt.Printf("🎉 [最終成功] A 餘額: %s, B 餘額: %s\n", resA.Kvs[0].Value, resB.Kvs[0].Value)
}
