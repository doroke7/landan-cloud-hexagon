package main

import (
	"fmt"
	"sync"
)

// Model  是單例 struct
type Model struct {
	name string
}

// 私有變量，存放唯一實例
var instance *Model
var once sync.Once

func NewModel(name string) *Model {
	return &Model{name: name}
}

//  once 直接寫在 method 裡面 變成 ，保證所有方法都只會做一次，這個在單例沒問題，

func Singletone(name string) *Model {

	// once 底層代碼就是先看有沒有初始化 再上鎖的思維/。
	// once 底層代碼就是先看有沒有初始化 再上鎖的思維/。
	// once 底層代碼就是先看有沒有初始化 再上鎖的思維/。
	// once 底層代碼就是先看有沒有初始化 再上鎖的思維/。
	// once 底層代碼就是先看有沒有初始化 再上鎖的思維/。
	// once 底層代碼就是先看有沒有初始化 再上鎖的思維/。

	once.Do(func() {
		instance = NewModel(name)
	})
	return instance
}

// 方法（模擬 getter/setter）
func (m *Model) GetName() string {
	return m.name
}

func (m *Model) SetName(name string) {
	m.name = name
}

func main() {
	oModel := Singletone("Tom")

	fmt.Println(oModel)
}
