package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
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

	// 初始化餘額：A = 300, B = 100
	_, _ = cli.Put(ctx, accA, "300")
	_, _ = cli.Put(ctx, accB, "100")

	// 🔄 因為 Txn 是底層原語、不內建重試，所以我們必須自己手動寫 for 迴圈來應對樂觀鎖衝突
	for {
		// 1. 【Read 階段】先撈出 A 和 B 目前的值與各自的 ModRevision (版本號)
		respA, _ := cli.Get(ctx, accA)
		respB, _ := cli.Get(ctx, accB)

		kvA := respA.Kvs[0]
		kvB := respB.Kvs[0]

		balanceA, _ := strconv.Atoi(string(kvA.Value))
		balanceB, _ := strconv.Atoi(string(kvB.Value))

		// 💡 記下這兩個 Key 此刻的版本防偽標籤
		revA := kvA.ModRevision
		revB := kvB.ModRevision

		fmt.Printf("👀 [Read] A餘額:%d(rev:%d), B餘額:%d(rev:%d)\n", balanceA, revA, balanceB, revB)

		if balanceA < 100 {
			fmt.Println("❌ 轉帳失敗：A 餘額不足")
			break
		}

		// 在 Go 程式裡計算好扣減後的新數值
		newBalanceA := balanceA - 100
		newBalanceB := balanceB + 100

		// 2. 【Compare & Write 階段】直接調用底層的 Txn 包裹一次送出
		txnResp, err := cli.Txn(ctx).
			// 🧱 If：檢查 A 和 B 的版本號是否跟剛才讀取時一模一樣
			If(
				clientv3.Compare(clientv3.ModRevision(accA), "=", revA),
				clientv3.Compare(clientv3.ModRevision(accB), "=", revB),
			).
			// 🧱 Then：條件完全成立，執行原子的寫入操作
			Then(
				clientv3.OpPut(accA, strconv.Itoa(newBalanceA)),
				clientv3.OpPut(accB, strconv.Itoa(newBalanceB)),
			).
			// 🧱 Else：條件不成立（被別人插隊了），這裡我們選擇什麼都不做（空操作）
			Else().
			// 🚀 正式提交給 etcd 叢集
			Commit()

		if err != nil {
			log.Fatalf("Txn 執行出錯: %v", err)
		}

		// 3. 【判斷結果】
		if txnResp.Succeeded {
			fmt.Println("🎉 [Success] etcd Txn 樂觀鎖校驗成功！轉帳完成。")
			break // 成功了，跳出 for 迴圈
		} else {
			// 💥 這就是樂觀鎖衝突！說明在 Read 到 Commit 之間，A 或 B 的版本號變了
			fmt.Println("💥 [Conflict] 有人插隊修改了帳戶！Txn 被拒絕，正在手動重試...")
			time.Sleep(50 * time.Millisecond) // 稍微錯開併發
		}
	}
}
