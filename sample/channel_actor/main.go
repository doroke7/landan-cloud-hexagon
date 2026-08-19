package main

import "fmt"

type Message struct {
	Type  string
	Value int
	Reply chan int
}

type CounterActor struct {
	mailbox chan Message
}

func NewCounterActor() *CounterActor {

	oCounterActor := &CounterActor{
		mailbox: make(chan Message),
	}

	go oCounterActor.run()

	return oCounterActor
}

func (oSelf *CounterActor) run() {

	iCount := 0

	for oMessage := range oSelf.mailbox {

		switch oMessage.Type {

		case "ADD":
			iCount += oMessage.Value

		case "GET":
			oMessage.Reply <- iCount
		}
	}
}

func (oSelf *CounterActor) Send(oMessage Message) {

	oSelf.mailbox <- oMessage
}

func main() {

	oCounter := NewCounterActor()

	oCounter.Send(Message{
		Type:  "ADD",
		Value: 10,
	})

	oChannel := make(chan int)

	oCounter.Send(Message{
		Type:  "GET",
		Reply: oChannel,
	})

	fmt.Println(<-oChannel)
}
