
# 由於 main.go 跟 user.go 同屬於 都是 package main。
# 編譯器在執行單一檔案時不會自動 import 同目錄的其他檔案。

# 解決辦法A: 另起一個目錄，用不同的package 名稱，再 import 近來            (是推薦)
# 解決辦法A: 同時執行兩個 go 文件。如 go run main.go user.go 或 go run . (不推薦)
# 解決辦法C: main 寫在同一個 go檔案                                    (不推薦)
