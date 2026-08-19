package main

import "fmt"

func main() {
	oBoxs := make(chan chan chan int, 1)

	oPapers := make(chan chan int, 1)

	oNumbers := make(chan int, 1)

	oBoxs <- oPapers

	oPapers <- oNumbers

	oNumbers <- 100

	for oBox := range oBoxs {
		oPapers := oBox

		for oPaper := range oPapers {
			oNumbers := oPaper

			for iNumber := range oNumbers {
				fmt.Println(iNumber)
			}

		}
	}
}
