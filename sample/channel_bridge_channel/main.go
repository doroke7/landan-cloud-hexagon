package main

import "fmt"

func bridge(oInputStream <-chan <-chan int) <-chan int {

	oOuputChannel := make(chan int)

	go func() {

		defer close(oOuputChannel)

		for oInputChannel := range oInputStream {

			for oInputValue := range oInputChannel {

				oOuputChannel <- oInputValue
			}
		}

	}()

	return oOuputChannel
}

func main() {

	oStream := make(chan (<-chan int))

	go func() {

		defer close(oStream)

		for i := 0; i < 3; i++ {

			oCh := make(chan int)

			go func(n int) {

				defer close(oCh)

				oCh <- n

				oCh <- n * 10

			}(i)

			oStream <- oCh
		}

	}()

	for oValue := range bridge(oStream) {

		fmt.Println(oValue)
	}
}
