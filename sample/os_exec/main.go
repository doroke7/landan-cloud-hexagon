package main

import (
	"fmt"
	"os/exec"
)

func main() {
	out, err := exec.Command("lscpu").Output() // Linux/macOS 可以改成 "sysctl -a | grep machdep.cpu" (mac)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
