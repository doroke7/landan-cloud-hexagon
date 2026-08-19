package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const maxChairs = 3 // 等候室只有 3 張椅子

type BarberShop struct {
	mu           sync.Mutex
	cond         *sync.Cond
	waitingCount int  // 目前在等候室坐著的顧客人數
	barberActive bool // 理髮師是否正在理髮（true: 忙碌, false: 睡覺/閒置）
}

func main() {
	shop := &BarberShop{}
	shop.cond = sync.NewCond(&shop.mu) // 將條件變數與鎖綁定

	// 1. 啟動理髮師執行緒
	go shop.barber()

	// 2. 模擬陸續進門的顧客
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(customerID int) {
			defer wg.Done()
			// 隨機間隔時間，模擬顧客隨機上門
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			shop.customer(customerID)
		}(i)
	}

	wg.Wait()
	time.Sleep(2 * time.Second) // 讓理髮師處理完最後的人
	fmt.Println("💈 今天的顧客都處理完了，下班！")
}

// 💇‍♂️ 理髮師的行為
func (shop *BarberShop) barber() {
	for {
		shop.mu.Lock()

		// 如果等候室沒人，理髮師就睡覺
		for shop.waitingCount == 0 {
			shop.barberActive = false
			fmt.Println("💤 店裡沒人，理髮師在理髮椅上呼呼大睡...")

			// shop.cond.Wait() 會釋放鎖，並阻塞在這裡，直到被顧客 Signal 喚醒。
			// 被喚醒時，它會重新自動自動獲取鎖（Lock）。

			//  來確保「檢查椅子數量」、「喚醒理髮師」這兩個動作是不可分割的原子操作
			//  來確保「檢查椅子數量」、「喚醒理髮師」這兩個動作是不可分割的原子操作
			//  來確保「檢查椅子數量」、「喚醒理髮師」這兩個動作是不可分割的原子操作
			//  來確保「檢查椅子數量」、「喚醒理髮師」這兩個動作是不可分割的原子操作
			//  來確保「檢查椅子數量」、「喚醒理髮師」這兩個動作是不可分割的原子操作
			shop.cond.Wait()
		}

		// 被喚醒或發現等候室有人，開始工作
		shop.barberActive = true
		shop.waitingCount-- // 叫等候室的第一個人進來理髮，釋放一張椅子
		fmt.Printf("✂️  理髮師被喚醒或叫號，開始幫一位顧客理髮（等候室剩餘椅子: %d）\n", maxChairs-shop.waitingCount)

		// 喚醒其他可能因為特殊狀態卡住的執行緒（在標準理髮師問題中，主要是通知剛進門的顧客理髮師醒了）
		shop.cond.Signal()
		shop.mu.Unlock()

		// 模擬理髮需要花費的時間（在外層解鎖後執行，不佔用鎖）
		time.Sleep(1 * time.Second)
	}
}

// 🚶 顧客的行為
func (shop *BarberShop) customer(id int) {
	shop.mu.Lock()
	fmt.Printf("🚶 顧客 %d 走進店裡...", id)

	// 如果等候室椅子滿了，直接離開
	if shop.waitingCount >= maxChairs {
		fmt.Printf("❌ 靠！等候室滿了，顧客 %d 直接掉頭離開。\n", id)

		//  來確保「檢查椅子數量」、「喚醒理髮師」這兩個動作是不可分割的原子操作
		//  來確保「檢查椅子數量」、「喚醒理髮師」這兩個動作是不可分割的原子操作
		//  來確保「檢查椅子數量」、「喚醒理髮師」這兩個動作是不可分割的原子操作
		//  來確保「檢查椅子數量」、「喚醒理髮師」這兩個動作是不可分割的原子操作
		//  來確保「檢查椅子數量」、「喚醒理髮師」這兩個動作是不可分割的原子操作
		shop.mu.Unlock()
		return
	}

	// 如果有空位，顧客留下來
	shop.waitingCount++
	fmt.Printf(" 坐下等候（等候室目前人數: %d/%d）\n", shop.waitingCount, maxChairs)

	// 發出信號：如果是因為沒人理髮師在睡覺，這個 Signal 會立刻把理髮師叫醒
	shop.cond.Signal()

	// 顧客必須等待，直到理髮師「不忙」且輪到自己（雖然本例簡化了顧客剪完的等待，主要聚焦在叫醒與排隊）
	for shop.barberActive && shop.waitingCount > 0 {
		shop.cond.Wait()
	}

	shop.mu.Unlock()
}
