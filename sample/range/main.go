package main

import "fmt"

func main() {
	oChannel1 := make(chan int, 10)
	oChannel2 := make(chan int, 10)

	// for-break 等價 for range
	for iNumber1 := range oChannel1 {
		fmt.Println(iNumber1)
	}

	// for-break 等價 for range
	for {
		iNumber2, bOk2 := <-oChannel2

		if !bOk2 {
			break
		}

		fmt.Println(iNumber2)

	}
}
