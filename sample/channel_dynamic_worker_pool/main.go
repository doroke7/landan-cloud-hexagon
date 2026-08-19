package main

import (
	"fmt"
	"time"
)

type Pool struct {
	jobs chan int
}

func NewPool() *Pool {

	return &Pool{
		jobs: make(chan int, 100),
	}
}

func (oSelf *Pool) AddWorker(id int) {

	go func() {

		for job := range oSelf.jobs {

			fmt.Println("worker", id, "process", job)
			time.Sleep(2 * time.Second)
		}

	}()

}

func main() {

	pool := NewPool()

	// 動態增加 worker

	pool.AddWorker(1)
	pool.AddWorker(2)
	pool.AddWorker(3)

	for i := 0; i < 10; i++ {

		pool.jobs <- i
	}

	time.Sleep(5 * time.Second)
}
