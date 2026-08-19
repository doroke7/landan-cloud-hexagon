package main

func main() {
	// map-bool   做成的 set 記憶體較多（1 bool   = 1byte 大小） ，
	// map-struct 做成的 set 記憶體較少（1 struct = 0byte 大小）
	// map set 建議用 map-struct
	
	oUsers := make(map[string]struct{})

	// 寫入：使用一個全域空結構體或直接初始化
	oUsers["11"] = struct{}{}

	// 判斷是否存在
	if _, ok := oUsers["11"]; ok {
		// 存在
	}
}
