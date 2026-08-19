package main

import "fmt"

// Producer
func producer() <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for i := 1; i <= 5; i++ {
			out <- i
		}
	}()

	return out
}

// 平方
func square(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for v := range in {
			out <- v * v
		}
	}()

	return out
}

// 乘二
func double(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for v := range in {
			out <- v * 2
		}
	}()

	return out
}

func main() {

	p := producer()

	s := square(p)

	d := double(s)

	for v := range d {
		fmt.Println(v)
	}
}
