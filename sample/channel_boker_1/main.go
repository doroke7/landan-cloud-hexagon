package main

import (
	"fmt"
	"sync"
	"time"
)

type Event struct {
	Message string
}

func NewBoker() *Broker {
	return &Broker{}
}

type Broker struct {
	subscribers []chan Event
	mutex       sync.Mutex
}

// channel： 一推送 一個message 只能被
func (oSelf *Broker) Subscribe() <-chan Event {

	oChannel := make(chan Event, 10)

	oSelf.mutex.Lock()
	defer oSelf.mutex.Unlock()

	oSelf.subscribers = append(
		oSelf.subscribers,
		oChannel,
	)

	return oChannel
}

// 發布
func (oSelf *Broker) Publish(oEvent Event) {

	oSelf.mutex.Lock()
	defer oSelf.mutex.Unlock()

	// 加鎖是 避免 訂閱突然增加人，但是同時 推送資料

	for _, oChannel := range oSelf.subscribers {

		select {

		case oChannel <- oEvent:
			// 發送成功

		default:
			// subscriber 太慢，丟棄
		}
	}
}

// 關閉所有訂閱
func (oSelf *Broker) Close() {

	oSelf.mutex.Lock()
	defer oSelf.mutex.Unlock()

	for _, oChannel := range oSelf.subscribers {
		close(oChannel)
	}
}

func main() {

	broker := NewBoker()

	// subscriber 1
	sub1 := broker.Subscribe()

	go func() {
		for event := range sub1 {
			fmt.Println(
				"sub1:",
				event.Message,
			)
		}
	}()

	// subscriber 2
	sub2 := broker.Subscribe()

	go func() {
		for event := range sub2 {
			fmt.Println(
				"sub2:",
				event.Message,
			)
		}
	}()

	// =-----------------------------------------------
	// publisher
	for i := 1; i <= 3; i++ {

		broker.Publish(Event{
			Message: fmt.Sprintf(
				"event %d",
				i,
			),
		})

		time.Sleep(time.Second)
	}

	broker.Close()

	time.Sleep(time.Second)
}
