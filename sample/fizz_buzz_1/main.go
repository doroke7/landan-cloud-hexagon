package main

import (
	"fmt"
	"sync"
)

/*
  Fizzbuzz問題的本質：
「一個全域的單一資料輸入來源（1 到 N），
分發給 N種 不同的消費隊列（Fizz、Buzz、FizzBuzz、Number），
但最卡手的是：這 1->2->3  ... -> N 個消費隊列之間，還必須維持絕對的『循序（Sequential）』。



*/

func main() {
	n := 15

	// 宣告 4 個通道，作為 4 個執行緒的接力棒
	chFizz := make(chan int)     // 負責 3 的倍數
	chBuzz := make(chan int)     // 負責 5 的倍數
	chFizzBuzz := make(chan int) // 負責 15 的倍數
	chNumber := make(chan int)   // 負責普通數字

	var wg sync.WaitGroup
	wg.Add(4)

	// 🧵 執行緒 A: Fizz
	go func() {
		defer wg.Done()
		for i := range chFizz {
			fmt.Println("Fizz")
			dispatch(i+1, n, chFizz, chBuzz, chFizzBuzz, chNumber) // 把接力棒交棒給下一個人
		}
	}()

	// 🧵 執行緒 B: Buzz
	go func() {
		defer wg.Done()
		for i := range chBuzz {
			fmt.Println("Buzz")
			dispatch(i+1, n, chFizz, chBuzz, chFizzBuzz, chNumber)
		}
	}()

	// 🧵 執行緒 C: FizzBuzz
	go func() {
		defer wg.Done()
		for i := range chFizzBuzz {
			fmt.Println("FizzBuzz")
			dispatch(i+1, n, chFizz, chBuzz, chFizzBuzz, chNumber)
		}
	}()

	// 🧵 執行緒 D: 普通數字
	go func() {
		defer wg.Done()
		for i := range chNumber {
			fmt.Println(i)
			dispatch(i+1, n, chFizz, chBuzz, chFizzBuzz, chNumber)
		}
	}()

	// 🏁 點火發射：先把數字 1 送給合適的通道（由調度函式決定）
	dispatch(1, n, chFizz, chBuzz, chFizzBuzz, chNumber)

	wg.Wait()
}

// 🎯 裁判/調度員：檢查目前數字該交給哪一個執行緒（保證 Check-then-Act 的絕對順序與原子性）
func dispatch(i, n int, chFizz, chBuzz, chFizzBuzz, chNumber chan int) {
	if i > n {
		// 超過範圍，關閉所有通道，讓所有 Goroutine 安全下班
		close(chFizz)
		close(chBuzz)
		close(chFizzBuzz)
		close(chNumber)
		return
	}

	// 根據精準的數學條件，把「目前的數字 i」當作接力棒塞給指定的通道
	if i%3 == 0 && i%5 == 0 {
		chFizzBuzz <- i
	} else if i%3 == 0 {
		chFizz <- i
	} else if i%5 == 0 {
		chBuzz <- i
	} else {
		chNumber <- i
	}
}
