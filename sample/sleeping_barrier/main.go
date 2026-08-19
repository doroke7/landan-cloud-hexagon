package main

import (
	"fmt"
	"time"
)

func barber(customers <-chan int) {

	for {

		customer := <-customers

		fmt.Println("barber cutting customer", customer)
		time.Sleep(2 * time.Second)
		fmt.Println("customer", customer, "done")
	}
}

func main() {

	// 3 張等待椅
	customers := make(chan int, 3)

	go barber(customers)

	for i := 1; i <= 10; i++ {

		select {

		case customers <- i:
			fmt.Println("customer", i, "sit waiting")

		default:

			fmt.Println("customer", i, "leave")

		}

		time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(10 * time.Second)
}
