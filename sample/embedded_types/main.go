package main

import "fmt"

type Base struct {
	Name string
}

func (b *Base) SayHello() {
	fmt.Println("Hello, I am", b.Name)
}

type Player struct {
	Base // 這裡就是 Embedded Types
}

func main() {
	u := Player{Base: Base{Name: "Gemini"}}
	// 魔法發生在這裡：User 本身沒定義 SayHello，
	// 但因為嵌入了 Base，SayHello 被「提升」到了 User 層級。
	u.SayHello()
}
