package main

import (
	"fmt"
	"time"
)

func worker(
	jobs <-chan int,
	stop <-chan struct{},
) {

	for {

		select {

		case job := <-jobs:

			fmt.Println(
				"process",
				job,
			)

			time.Sleep(time.Second)

		case <-stop:

			fmt.Println(
				"stop signal received",
			)

			// drain queue
			for job := range jobs {

				fmt.Println(
					"drain process",
					job,
				)
			}

			return
		}
	}
}

func main() {

	jobs := make(chan int, 10)

	stop := make(chan struct{})

	go worker(
		jobs,
		stop,
	)

	for i := 1; i <= 5; i++ {

		jobs <- i
	}

	time.Sleep(2 * time.Second)

	close(stop)

	close(jobs)

	time.Sleep(5 * time.Second)
}
