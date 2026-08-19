package main

import "fmt"

func main() {

	cAdd := func(iX int, iY int) int {
		return iX + iY
	}

	iResult := cAdd(1, 2)

	fmt.Printf("Result: %d\n", iResult)
}
