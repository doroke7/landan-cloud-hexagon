package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	u1 := new(User)
	u1.Name = "Gemini"

	u2 := &User{
		Name: "Gemini2",
	}
	fmt.Printf("%T\n", u1) // 輸出: *main.User (是一個指針)
	fmt.Printf("%T\n", u2) // 輸出: *main.User (是一個指針)

}

/**
特性	new(T)	                                    make(T, args)
適用類型	所有類型 (常用於結構體)。	                僅限 Slice, Map, Channel。
返回內容	指針 (*T)。	                            該類型本身 (Value)。
初始化行為	僅清空記憶體（零值），不進行複雜初始化。	會初始化內部的數據結構（如設定切片長度、哈希表空間）。
*/
