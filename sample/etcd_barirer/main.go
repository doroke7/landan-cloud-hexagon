// 🧪 【單一程式、雙角色測試指南】
// 1. 先把此檔案儲存為 main.go
// 2. 啟動一個或多個 Worker（它們會卡住排隊）：
//    go run main.go worker
// 3. 打開另一個終端機，啟動 Master（它會負責開門）：
//    go run main.go master

package main

import (
	"fmt"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

func main() {
	// 判斷終端機有沒有帶參數
	if len(os.Args) < 2 {
		fmt.Println("❌ 請指定角色！用法：go run main.go [master|worker]")
		return
	}
	role := os.Args[1]

	cli, _ := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	defer cli.Close()

	// 認準兩道門
	barrier1 := recipe.NewBarrier(cli, "/step1_gate")
	barrier2 := recipe.NewBarrier(cli, "/step2_gate")

	// =========================================================================
	// 👑 角色 A：如果你是 MASTER (工頭)
	// =========================================================================
	if role == "master" {
		fmt.Println("🎬 [Master Mode] 啟動...")

		_ = barrier1.Hold() // IMPORTANT: Master 專有，鎖上第一道大門
		_ = barrier2.Hold() // IMPORTANT: Master 專有，鎖上第二道大門
		fmt.Println("🔒 大門已全部鎖死，請啟動 worker 開始排隊。5 秒後開啟第一道門...")

		time.Sleep(5 * time.Second)
		fmt.Println("📢 開啟 Step 1 大門！")
		_ = barrier1.Release() // IMPORTANT: Master 專有，拆除第一道門

		time.Sleep(5 * time.Second)
		fmt.Println("📢 開啟 Step 2 大門！")
		_ = barrier2.Release() // IMPORTANT: Master 專有，拆除第二道門

		fmt.Println("✨ Master 任務結束，下班！")
		return
	}

	// =========================================================================
	// 🏃 角色 B：如果你是 WORKER (苦力)
	// =========================================================================
	if role == "worker" {
		nodeName := fmt.Sprintf("Worker-PID-%d", os.Getpid())
		fmt.Printf("🎬 [%s] 啟動...\n", nodeName)

		fmt.Printf("🛑 [%s] 來到 [Step 1 閘門] 前... 卡住等待...\n", nodeName)
		_ = barrier1.Wait() // IMPORTANT: Worker 只能被動等待

		fmt.Printf("🚀 [%s] Step 1 通過！執行第一階段工作中...\n", nodeName)
		time.Sleep(2 * time.Second)

		fmt.Printf("🛑 [%s] 第一階段做完，來到 [Step 2 閘門] 前... 再次卡住...\n", nodeName)
		_ = barrier2.Wait() // IMPORTANT: Worker 再次被動等待

		fmt.Printf("👑 [%s] 兩階段大功告成！\n", nodeName)
		return
	}

	fmt.Println("❌ 未知的角色，請輸入 master 或 worker")
}
