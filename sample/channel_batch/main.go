package main

import (
	"fmt"
	"time"
)

func batchWorker(ch <-chan int) {

	batch := make([]int, 0, 5)

	for value := range ch {

		batch = append(batch, value)

		if len(batch) == 5 {

			process(batch)

			batch = batch[:0]
		}
	}
}

func process(batch []int) {

	fmt.Println(
		"batch:",
		batch,
	)

	time.Sleep(time.Second)
}

func main() {

	ch := make(chan int)

	go batchWorker(ch)

	for i := 1; i <= 12; i++ {

		ch <- i
	}

	close(ch)

	time.Sleep(time.Second)
}
