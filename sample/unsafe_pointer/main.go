package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// 嚴謹說，其實 golang 本身沒有私有屬性。 小寫私有完全是編譯控制的
// 利用反射可以 取得小寫屬性

type Student struct {
	Name     string
	password string // 小寫私有，正常情況下外部無法存取
}

func main() {
	oStudent := Student{Name: "Gemini", password: "Secret123"}
	fmt.Println("1. 原始資料:", oStudent)

	// --- 黑魔法開始 ---

	// 1. 使用反射 (reflect) 取得結構體資訊
	// 我們需要知道 password 這個欄位在記憶體裡的「偏移量 (Offset)」
	oVal := reflect.ValueOf(&oStudent).Elem()
	oField, _ := oVal.Type().FieldByName("password")
	offset := oField.Offset

	// 2. 取得結構體的起始位址
	pAddr := unsafe.Pointer(&oStudent)

	// 3. 起始位址 + 偏移量 = password 的真正記憶體座標
	// uintptr(structAddr) 將指標轉成數字，加上偏移量後，再轉回 unsafe.Pointer
	pPasswordAddr := unsafe.Pointer(uintptr(pAddr) + offset)

	// 4. 強行轉換並讀取
	// 我們知道 pPasswordAddr 是無狀態指標，改成 string指標 即 *string
	ptrToPassword := (*string)(pPasswordAddr)

	fmt.Println("2. 成功讀取私有屬性:", *ptrToPassword)

	// 5. 甚至可以強行修改它
	*ptrToPassword = "HackedByUnsafe"
	fmt.Println("3. 修改後的資料:", oStudent)
}
