package bootstrap

import "fmt"

// gormWriter 是 gorm 的 logger.New 需要的 logger.Writer —— 只要有 Printf(string, ...any) 就行。
// 這裡用 fmt 直接印到 stdout，取代原本 log.New(os.Stdout, "\r\n", log.LstdFlags) 那顆 stdlib logger
// （gorm 的 trace 字串沒有結尾換行，靠 Writer 自己補，所以這裡 format 後面接一個 \n）。
type gormWriter struct{}

func (gormWriter) Printf(sFormat string, aArgs ...any) {
	fmt.Printf(sFormat+"\n", aArgs...)
}
