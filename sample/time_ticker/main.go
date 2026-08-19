package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		t := <-ticker.C
		fmt.Println("收到 Tick:", t.Format("15:04:05"))
	}
}
