package main

import (
	"fmt"
	"sync"
)

func main() {
	var mutex sync.Mutex

	fmt.Println("1. 準備第一次加鎖...")
	mutex.Lock() // 成功拿到鎖
	fmt.Println("   成功拿到第一次鎖！")

	fmt.Println("2. 準備第二次加鎖（注意：此時還沒解鎖喔！）...")
	
	// 💥 這裡會直接卡死（死鎖）！
	// 因為同一個 Goroutine 試圖去搶一把「已經被自己鎖住、且還沒放開」的鎖。
	mutex.Lock() 

	fmt.Println("3. 這行永遠不會被印出來...")
	mutex.Unlock()
	mutex.Unlock()
}