package main

import (
	"fmt"
)

func main() {
	//  0  1  2  3  4  5  6  7  8  9
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	s1 := slice[2:5] // [2 3 4] len=3 cap=8

	// s1 的底層:
	// index: 0 1 2 3 4 5 6 7
	// value: 2 3 4 5 6 7 8 9

	s2 := s1[2:6:7] // [4 5 6 7] len=4 cap=5

	// s2 的底層:
	// index: 0 1 2 3 4
	// value: 4 5 6 7 8

	s2 = append(s2, 100) // [4 5 6 7 100]，直接覆蓋 slice[8]

	s2 = append(s2, 200) // cap 不夠，重新配置新陣列

	s1[2] = 20 // 修改 slice[4]

	fmt.Println(s1) // [2 3 20]

	fmt.Println(s2) // [4 5 6 7 100 200]

	fmt.Println(slice) // [0 1 2 3 20 5 6 7 100 9]
}
