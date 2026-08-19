package main

import "fmt"

func main() {

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		ch1 <- 100
	}()

	go func() {
		ch2 <- 200
	}()

	for i := 0; i < 2; i++ {

		select {

		case v := <-ch1:
			fmt.Println("ch1:", v)

			// disable ch1
			ch1 = nil

		case v := <-ch2:
			fmt.Println("ch2:", v)

			// disable ch2
			ch2 = nil
		}
	}
}
