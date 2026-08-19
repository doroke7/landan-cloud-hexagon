package main

import (
	"fmt"
	"time"
)

type Barrier struct {
	done chan struct{}
}

func NewBarrier(count int) *Barrier {

	b := &Barrier{
		done: make(chan struct{}),
	}

	go func() {

		for i := 0; i < count; i++ {

			<-b.done
		}

		close(b.done)

	}()

	return b
}

func (b *Barrier) Wait() {

	// 通知自己到達
	b.done <- struct{}{}

	// 等待所有人完成
	<-b.done
}

func worker(id int,
	barrier *Barrier,
) {

	fmt.Println(
		"worker",
		id,
		"before barrier",
	)

	time.Sleep(
		time.Duration(id) * time.Second,
	)

	fmt.Println(
		"worker",
		id,
		"arrived",
	)

	barrier.Wait()

	fmt.Println(
		"worker",
		id,
		"after barrier",
	)
}

func main() {

	barrier := NewBarrier(3)

	for i := 1; i <= 3; i++ {

		go worker(
			i,
			barrier,
		)

	}

	time.Sleep(6 * time.Second)
}
