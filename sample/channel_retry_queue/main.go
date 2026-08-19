package main

import (
	"fmt"
	"time"
)

type Job struct {
	ID int
}

func worker(
	jobs chan Job,
	retry chan Job,
) {

	for {

		select {

		case job := <-jobs:

			fmt.Println(
				"process job",
				job.ID,
			)

			// 模擬失敗
			if job.ID%2 == 0 {

				fmt.Println(
					"failed",
					job.ID,
				)

				// 放入 retry queue
				go func() {

					time.Sleep(2 * time.Second)

					retry <- job

				}()

				continue
			}

			fmt.Println(
				"success",
				job.ID,
			)

		case oRetry := <-retry:

			fmt.Println(
				"retry job",
				oRetry.ID,
			)

			jobs <- oRetry
		}
	}
}

func main() {

	jobs := make(chan Job, 10)

	retry := make(chan Job, 10)

	go worker(
		jobs,
		retry,
	)

	for i := 1; i <= 10; i++ {

		jobs <- Job{
			ID: i,
		}
	}

	time.Sleep(10 * time.Second)
}
