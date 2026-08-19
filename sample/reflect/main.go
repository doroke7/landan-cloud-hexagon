package main

import (
	"fmt"
	"reflect"
)

type Config struct {
	Directory string
	Size      string
	Number    string
}

func main() {
	oConfig := Config{
		Directory: "/usr/local/logs",
		Size:      "10M",
		Number:    "100",
	}

	// 取得反射值物件
	oValue := reflect.ValueOf(oConfig)

	fmt.Println("類型:", oValue.Type()) // main.Config
	fmt.Println("種類:", oValue.Kind()) // struct

	// 動態獲取欄位內容
	f := oValue.FieldByName("Directory")
	fmt.Println("內容:", f.String()) // /usr/local/logs

	fmt.Println("mapping oConfig[\"Number\"] 內容:", mapping(oConfig, "Number")) // /usr/local/logs

}

/*
*

		Go 語言 這種強型別的靜態語言裡，反射 (Reflection) 最核心的用途之一，
		就是把原本「寫死」的Struct欄位，變成像 Map 一樣可以用 「字串（文字索引）」 來存取的對象。
	    如果你不用反射，這樣的代碼有會變長非常多的 if-else 取值
*/
func mapping(oConfig Config, sName string) string {
	oValue := reflect.ValueOf(oConfig)
	oField := oValue.FieldByName(sName)

	return oField.String()
}
