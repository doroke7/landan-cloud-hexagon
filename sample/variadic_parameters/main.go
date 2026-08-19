package main

// JS 的 ... 是作用在變量上的 代表 變量的收集或展開
// Go 的 ... 是作用在型態上的 代表 他這個變量型別是收集的型態

func Add(aArgs ...int) int {
	var iResult int = 0
	for _, iArg := range aArgs {
		iResult = iArg + iResult
	}
	return iResult
}

func main() {
	Add(1, 2, 3, 4, 5, 7, 8)

	aNumbers := []int{11, 12, 13, 14, 15, 16, 17, 18, 19}

	Add(aNumbers...)
}
