package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// 1. 打開檔案
	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatalf("無法開啟檔案: %v", err)
	}
	// 記得在程式結束時關閉檔案，釋放系統資源
	defer file.Close()

	// 2. 建立一個專門用來「掃描」該檔案的 Scanner
	// 它會在記憶體裡開闢快取，有效率地批次讀取
	scanner := bufio.NewScanner(file)

	fmt.Println("====== 開始逐行讀取檔案 ======")

	// 3. 使用 for scanner.Scan() 開始循環
	// 每次呼叫 Scan()，它就會自動讀到下一個換行符號（\n），並回傳 true
	// 如果讀到檔案末尾（EOF）或出錯，它就會回傳 false，迴圈就會自動停止
	for scanner.Scan() {
		// scanner.Text() 會直接回傳剛剛讀到的那一整行字串（而且已經幫你把結尾的 \n 刪掉了！）
		line := scanner.Text()

		fmt.Printf("讀取到一行數據: %s\n", line)
		time.Sleep(1 * time.Second)
	}

	// 4. 檢查迴圈結束是因為「正常讀完」還是「中途出錯」
	if err := scanner.Err(); err != nil {
		log.Fatalf("讀取檔案過程中出錯: %v", err)
	}

	fmt.Println("====== 檔案讀取完畢 ======")
}
