package main

func main() {
	// IDE 會從 go.mod 查看有沒有，如果 go.mod 有指定代碼，但是對應的host 目錄沒有，他會偷偷下載
	// 所以 IDE 的流程是， 查看 go.mod 去標記的目錄找代碼（雖然是 虛擬機的目錄, 宿主機肯定沒有），發現沒有就偷偷下載一份
}
