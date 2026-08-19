package main

import (
	"fmt"
	"time"
)

type UnboundedQueue struct {
	in  chan int
	out chan int
}

func NewUnboundedQueue() *UnboundedQueue {

	oUnboundedQueue := &UnboundedQueue{
		in:  make(chan int),
		out: make(chan int),
	}

	go oUnboundedQueue.run()

	return oUnboundedQueue
}

func (oUnboundedQueue *UnboundedQueue) run() {

	aQueues := []int{}

	for {
		if len(aQueues) == 0 {

			value := <-oUnboundedQueue.in

			aQueues = append(aQueues, value)

			continue
		}

		select {

		case value := <-oUnboundedQueue.in:
			aQueues = append(aQueues, value)

		case oUnboundedQueue.out <- aQueues[0]:
			aQueues = aQueues[1:]
		}
	}
}

func (oUnboundedQueue *UnboundedQueue) Push(value int) {

	oUnboundedQueue.in <- value
}

func (oUnboundedQueue *UnboundedQueue) Pop() int {

	return <-oUnboundedQueue.out
}

func main() {

	oUnboundedQueue := NewUnboundedQueue()

	go func() {

		for i := 0; i < 10; i++ {

			oUnboundedQueue.Push(i)

			fmt.Println("push", i)
		}

	}()

	go func() {

		for {

			iValue := oUnboundedQueue.Pop()

			fmt.Println("pop", iValue)

			time.Sleep(time.Second)
		}

	}()

	time.Sleep(5 * time.Second)
}
