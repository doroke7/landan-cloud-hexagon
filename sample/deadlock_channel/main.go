package main

import (
	"fmt"
	"sync"
)

/*

範例二：Channel 互相等待（你等我送，我等你收）
在 Go 語言中，不只是 sync.Mutex 會死鎖，無緩衝的 Channel（Unbuffered Channel） 如果發送與接收的節奏沒對上，也超級容易製造死鎖。
這個例子裡，兩個 Goroutine 都想「先從對方那裡拿到資料，才肯把自己的資料送出去」。


*/

func main() {
	// 建立兩個無緩衝的 Channel
	ch1 := make(chan int)
	ch2 := make(chan int)
	
	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1
	go func() {
		defer wg.Done()
		fmt.Println("G1：我想先從 ch1 讀取資料...")
		
		data := <-ch1 // 💥 卡死！G1 在這裡停下來，等待有人放資料進 ch1
		
		fmt.Println("G1：讀到了！準備把資料送進 ch2...")
		ch2 <- data
	}()

	// Goroutine 2
	go func() {
		defer wg.Done()
		fmt.Println("G2：我想先從 ch2 讀取資料...")
		
		data := <-ch2 // 💥 卡死！G2 也在這裡停下來，等待有人放資料進 ch2
		
		fmt.Println("G2：讀到了！準備把資料送進 ch1...")
		ch1 <- data
	}()

	wg.Wait()
}