package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Model 是單例 struct
type Model struct {
	name string
}

// 私有變量，存放唯一實例
var instance *Model
var mu sync.Mutex

// initialized 用原子操作來標記是否已經初始化完成（0: 未完成, 1: 已完成）
// 使用原子操作（atomic）可以防止多個 CPU 核心之間發生「指令重排」或快取不一致的問題
var initialized uint32

func NewModel(name string) *Model {
	return &Model{name: name}
}

// once 就是先看有沒有初始化 再上鎖的思維。 避免無腦式上鎖，造成性能損失

func Singletone(name string) *Model {
	if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}

	// 沒初始化，大家開始排隊搶鎖
	mu.Lock()
	defer mu.Unlock()

	if initialized == 0 {
		instance = NewModel(name)

		atomic.StoreUint32(&initialized, 1)
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
