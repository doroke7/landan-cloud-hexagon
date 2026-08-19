package main

func main() {
	ch := make(chan int, 1)

	close(ch) // 關閉 channel

	ch <- 100 // ❌ Panic: send on closed channel
}
