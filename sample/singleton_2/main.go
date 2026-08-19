package main

import (
	"fmt"
	"sync"
)

// Model 是單例 struct
type Model struct {
	name string
}

// 私有變量，存放唯一實例
var instance *Model
var mu sync.Mutex // 只需要一把互斥鎖，不需要額外的 atomic 標記了

func NewModel(name string) *Model {
	return &Model{name: name}
}

// Singletone 用手動鎖（只檢查一次）實現單例
func Singletone(name string) *Model {
	// 不管三七二十一，進來的人全部排隊搶鎖
	mu.Lock()
	defer mu.Unlock()

	// 【只檢查一次】：因為鎖在最外層，進到這裡的人絕對是循序一對一執行的
	if instance == nil {
		instance = NewModel(name)
	}

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
	fmt.Println(oModel.GetName()) // 輸出: Tom
}
