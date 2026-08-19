package main

import (
	"fmt"
	"sync"
)

type Result struct {
	Value  int
	Square int
}

func worker(id int, oJobs <-chan int, oResults chan<- Result, oWg *sync.WaitGroup) {

	defer oWg.Done()

	for iJob := range oJobs {

		fmt.Println("worker", id, "process", iJob)

		oResults <- Result{
			Value:  iJob,
			Square: iJob * iJob,
		}
	}
}

func main() {

	oJobs := make(chan int)

	oResults := make(chan Result)

	var oWg sync.WaitGroup

	// 發送工作
	go func() {

		for i := 1; i <= 10; i++ {
			oJobs <- i
		}

		close(oJobs)

	}()

	go func() {

		oWg.Wait()

		close(oResults)

	}()

	// Scatter
	for iIndex := 1; iIndex <= 4; iIndex++ {

		oWg.Add(1)

		go worker(iIndex, oJobs, oResults, &oWg)
	}

	// Gather
	for iResult := range oResults {

		fmt.Println(
			"result:",
			iResult,
		)
	}
}
