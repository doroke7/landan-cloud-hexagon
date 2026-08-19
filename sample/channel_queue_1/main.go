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

		var iFirst int
		var oOutChannel chan int

		if len(aQueues) > 0 {
			iFirst = aQueues[0]
			oOutChannel = oUnboundedQueue.out

			/*
				這裡就是 Go select 裡一個很經典的技巧：
				利用 nil channel 永遠阻塞的特性，動態開關 select case。
				你的理解正確。

			*/
		}

		select {

		case value := <-oUnboundedQueue.in:

			aQueues = append(aQueues, value)

		case oOutChannel <- iFirst:

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
