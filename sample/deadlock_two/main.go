package main

import (
	"fmt"
	"sync"
	"time"
)

/**

範例一：資源互卡（你拿 A 等 B，我拿 B 等 A）
想像一下：小明手裡拿著畫筆，正在等小華把畫紙給他；
同一時間，小華手裡拿著畫紙，正在等小明把畫筆給他。
兩個人都堅持不放手，於是就永遠卡住了。


*/

func main() {
	// 用真實道具當作鎖，直覺多了！
	var lockPen sync.Mutex   // 筆的鎖
	var lockPaper sync.Mutex // 紙的鎖
	
	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1：小明
	go func() {
		defer wg.Done()
		
		fmt.Println("小明：搶到了【筆】🔒")
		lockPen.Lock() 

		// 故意睡一下，讓小華有時間去搶紙
		time.Sleep(100 * time.Millisecond)

		fmt.Println("小明：想拿【紙】...（等小華放開）")
		lockPaper.Lock() // 💥 卡死在這裡！因為紙被小華鎖住了
		
		fmt.Println("小明：太棒了，筆跟紙都拿到了，開始畫畫！")
		lockPaper.Unlock()
		lockPen.Unlock()
	}()

	// Goroutine 2：小華
	go func() {
		defer wg.Done()
		
		fmt.Println("小華：搶到了【紙】🔒")
		lockPaper.Lock() 

		// 故意睡一下，讓小明有時間去搶筆
		time.Sleep(100 * time.Millisecond)

		fmt.Println("小華：想拿【筆】...（等小明放開）")
		lockPen.Lock() // 💥 卡死在這裡！因為筆被小明鎖住了
		
		fmt.Println("小華：太棒了，紙跟筆都拿到了，開始畫畫！")
		lockPen.Unlock()
		lockPaper.Unlock()
	}()

	wg.Wait()
}