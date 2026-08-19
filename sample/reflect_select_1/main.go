package main

import (
	"fmt"
	"reflect"
)


/*
	reflect.Select 的用處：
	做一個動態case欄目的 select ，可以收任意 多個通道 （注意是通道多個的意思，不是元素）

*/
func main() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch3 := make(chan int, 1)
	ch4 := make(chan int, 1)

	ch1 <- 100
	ch2 <- 200
	ch3 <- 300
	ch4 <- 400

	close(ch1)
	close(ch2)
	close(ch3)
	close(ch4)

	// 動態管理 channel；之後增加 ch5，只要加到這裡即可。
	channels := []<-chan int{ch1, ch2, ch3, ch4}
	names := []string{"ch1", "ch2", "ch3", "ch4"}

	// 每個 channel 對應一個 SelectCase。
	cases := make([]reflect.SelectCase, len(channels))
	for i, ch := range channels {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}

	active := len(cases)

	// 直到所有 channel 都被關閉、停用。
	for active > 0 {
		// 這邊拉平了數據， 從2維度（4 個 channel k 個元素）變成一個維度
		chosen, recv, recvOK := reflect.Select(cases)

		if !recvOK {
			fmt.Printf("【%s】已關閉\n", names[chosen])

			// 關閉的 channel 永遠「可讀」。
			// 停用它，避免下一輪 reflect.Select 又一直選到它。
			cases[chosen].Chan = reflect.Value{}
			active--
			continue
		}

		// chosen 是動態的 channel index，不需要 switch。
		fmt.Printf(
			"【%s】收到資料: %d（型態: %s）\n",
			names[chosen],
			recv.Int(),
			recv.Type(),
		)
	}

	fmt.Println("所有 channel 都已處理完畢")
}