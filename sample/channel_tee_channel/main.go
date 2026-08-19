package main

import "fmt"

func tee(in <-chan int) (<-chan int, <-chan int) {

	out1 := make(chan int)
	out2 := make(chan int)

	go func() {

		defer close(out1)
		defer close(out2)

		for value := range in {

			// 複製給 out1
			select {
			case out1 <- value:
			}

			// 複製給 out2
			select {
			case out2 <- value:
			}
		}

	}()

	return out1, out2
}

func main() {

	input := make(chan int)

	go func() {

		defer close(input)

		for i := 1; i <= 3; i++ {
			input <- i
		}

	}()

	a, b := tee(input)

	for v := range a {
		fmt.Println("A:", v)
	}

	for v := range b {
		fmt.Println("B:", v)
	}
}
