package main

import "fmt"

// 1. 定義介面（鴨子合約）：只要會 Quack 的，就是鴨子！
type Duck interface {
	Quack() string
}

// ==========================================
// 2. 結構體 A：真正的生物鴨子
type RealDuck struct {
	Name string
}

// 真正的鴨子實現了 Quack 方法
func (r RealDuck) Quack() string {
	return "呱呱！我是真正的野生鴨子：" + r.Name
}

// ==========================================
// 3. 結構體 B：高科技機械鳥
// 💡 注意：這個結構體在定義時，完全沒有提到 "Duck" 這個字！
type RobotBird struct {
	SerialNumber string
}

// 機械鳥也實現了 Quack 方法（簽章跟 Duck 一模一樣）
func (b RobotBird) Quack() string {
	return "嗶嗶！機械核心啟動！發出電子音：呱呱！編號：" + b.SerialNumber
}

// ==========================================
// 4. 業務邏輯：只要符合 Duck 介面的物件，都能丟進這個函式
func LetTheDuckSpeak(d Duck) {
	fmt.Println(d.Quack())
}

func main() {
	// 宣告一隻真鴨子
	donald := RealDuck{Name: "唐老鴨"}

	// 宣告一隻機械鳥
	t800 := RobotBird{SerialNumber: "N95-G00"}

	fmt.Println("--- 執行鴨子型別測試 ---")

	// 測試真鴨子：走路像鴨子，叫聲像鴨子，所以放行！
	LetTheDuckSpeak(donald)

	// 測試機械鳥：雖然它是機器做的，但叫聲也像鴨子，Go 編譯器一樣放行！
	LetTheDuckSpeak(t800)

	fmt.Println("-----------------------")
}

/**

看起來，對於小型的輸入接口， go 直接用語言特性 鴨子型態 就能解決了 未必需要 adaptor pattern


*/
