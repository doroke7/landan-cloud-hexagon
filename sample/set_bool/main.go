package main

func main() {
	// map-bool   做成的 set 記憶體較多（1 bool   = 1byte 大小） ，
	// map-struct 做成的 set 記憶體較少（1 struct = 0byte 大小）
	// map set 建議用 map-struct

	oUsers := make(map[string]bool)
	oUsers["11"] = true

	if oUsers["11"] {
		// 存在
	}
}
