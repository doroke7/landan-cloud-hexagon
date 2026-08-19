package main

import (
	"fmt"
	"time"
)

func main() {

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {

		time.Sleep(time.Second)

		ch1 <- "database result"

	}()

	go func() {

		time.Sleep(2 * time.Second)

		ch2 <- "cache result"

	}()

	select {

	case v := <-ch1:

		fmt.Println(v)

	case v := <-ch2:

		fmt.Println(v)

	}
}
