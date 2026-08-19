package main

import (
	"fmt"
)

func main() {
	oChannel1 := make(chan int, 1)

	oChannel2 := make(chan int, 1)

	oChannel3 := make(chan int, 1)

	oChannel3 <- 11

	x := <-oChannel3

	oChannel2 <- x

	y := <-oChannel2

	oChannel1 <- y

	z := <-oChannel1
	fmt.Println("z=", z)

}
