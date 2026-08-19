package main

import (
	"fmt"
	"sync"
	"time"
)

type Event struct {
	Message string
}

type Broker struct {
	subscribers map[string]chan Event
	mutex       sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string]chan Event),
	}
}

// Subscribe
func (b *Broker) Subscribe(id string) <-chan Event {

	ch := make(chan Event, 10)

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.subscribers[id] = ch

	return ch
}

// Unsubscribe
func (b *Broker) Unsubscribe(id string) {

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if ch, ok := b.subscribers[id]; ok {

		close(ch)

		delete(
			b.subscribers,
			id,
		)
	}
}

// Publish
func (b *Broker) Publish(event Event) {

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	for _, ch := range b.subscribers {

		select {

		case ch <- event:
			// 發送

		default:
			// subscriber 太慢
		}
	}
}

func main() {

	broker := NewBroker()

	// subscriber A
	a := broker.Subscribe("email")

	go func() {

		for event := range a {

			fmt.Println(
				"email:",
				event.Message,
			)
		}

	}()

	// subscriber B
	b := broker.Subscribe("log")

	go func() {

		for event := range b {

			fmt.Println(
				"log:",
				event.Message,
			)
		}

	}()

	broker.Publish(Event{
		Message: "User Created",
	})

	time.Sleep(time.Second)

	broker.Unsubscribe("email")

	broker.Publish(Event{
		Message: "User Deleted",
	})

	time.Sleep(time.Second)
}
