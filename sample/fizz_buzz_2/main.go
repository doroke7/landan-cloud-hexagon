package main

import (
	"fmt"
	"sync"
)

// 這邊版本 不是串行

func main() {
	n := 15

	// 宣告 4 個通道，單純用來接收指揮官傳來的列印任務
	chFizz := make(chan int)
	chBuzz := make(chan int)
	chFizzBuzz := make(chan int)
	chNumber := make(chan int)

	var wg sync.WaitGroup
	wg.Add(4)

	// 🧵 小兵 A: 專門列印 Fizz
	go func() {
		defer wg.Done()
		for range chFizz { // 只要通道有東西就印，不用管數字是幾
			fmt.Println("Fizz")
		}
	}()

	// 🧵 小兵 B: 專門列印 Buzz
	go func() {
		defer wg.Done()
		for range chBuzz {
			fmt.Println("Buzz")
		}
	}()

	// 🧵 小兵 C: 專門列印 FizzBuzz
	go func() {
		defer wg.Done()
		for range chFizzBuzz {
			fmt.Println("FizzBuzz")
		}
	}()

	// 🧵 小兵 D: 專門列印普通數字
	go func() {
		defer wg.Done()
		for num := range chNumber {
			fmt.Println(num) // 只有他需要知道具體的數字
		}
	}()

	// 🎯 指揮官（Single Source of Truth）：負責所有的判斷與投餵
	for i := 1; i <= n; i++ {
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

	// 投餵結束，指揮官關閉所有通道，通知小兵們可以下班了
	close(chFizz)
	close(chBuzz)
	close(chFizzBuzz)
	close(chNumber)

	wg.Wait()
}
