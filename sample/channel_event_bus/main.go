package main

import (
	"fmt"
	"sync"
)

type Event struct {
	Name string
}

type EventBus struct {
	mu sync.Mutex

	subscribers []chan Event
}

func (b *EventBus) Subscribe() <-chan Event {

	ch := make(chan Event, 10)

	b.mu.Lock()

	b.subscribers = append(
		b.subscribers,
		ch,
	)

	b.mu.Unlock()

	return ch
}

func (b *EventBus) Publish(event Event) {

	b.mu.Lock()

	defer b.mu.Unlock()

	for _, ch := range b.subscribers {

		ch <- event
	}
}

func main() {

	bus := &EventBus{}

	sub1 := bus.Subscribe()

	sub2 := bus.Subscribe()

	go func() {

		for event := range sub1 {

			fmt.Println(
				"sub1:",
				event.Name,
			)
		}

	}()

	go func() {

		for event := range sub2 {

			fmt.Println(
				"sub2:",
				event.Name,
			)
		}

	}()

	bus.Publish(
		Event{
			Name: "OrderCreated",
		},
	)

}
