package main

func main() {

	ch := make(chan int, 4)

	go func() {

		ch <- 1
		ch <- 2
		ch <- 3
		ch <- 4
		ch <- 5

	}()

	for i := 1; i <= 5; i++ {

		println(<-ch)

		/*
			   channel 的 輸入 輸出 必須相當 相當


			// 前面 5 個輸入 channel
			// 後面 就同時要有 5 個輸出

		*/
	}
}
