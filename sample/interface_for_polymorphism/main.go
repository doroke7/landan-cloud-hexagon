package main

import "fmt"

/**
本身實現多形是 靠 interface
*/
// 定義介面
type VideoProcessor interface {
	Process() string
}

// 結構體 A：本地影片
type LocalVideo struct {
	FileName string
}

// 實作介面：使用「指標接收者 (*LocalVideo)」
// 這樣做可以修改內部狀態，且在大物件傳遞時效能更好
func (v *LocalVideo) Process() string {
	return fmt.Sprintf("正在處理本地檔案: %s", v.FileName)
}

// 結構體 B：串流影片
type StreamVideo struct {
	URL string
}

// 實作介面：使用「值接收者 (StreamVideo)」
func (v StreamVideo) Process() string {
	return fmt.Sprintf("正在處理串流: %s", v.URL)
}

type Test struct {
}

func (oSelf *Test) Process() string {
	return fmt.Sprintf("TEST 拉")

}

func RunTask(p VideoProcessor) {
	fmt.Println(p.Process())
}

func main() {
	lv := &LocalVideo{FileName: "movie.mp4"}
	sv := StreamVideo{URL: "https://live.com/test"}
	oTest := &Test{}

	// 多型展現：RunTask 接受任何實作了 VideoProcessor 的型別
	RunTask(lv)    // 傳入指標
	RunTask(sv)    // 傳入值
	RunTask(oTest) // 傳入值

}
