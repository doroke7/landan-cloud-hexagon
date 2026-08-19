package helper

import (
	_ "fmt"
)

/*
*
你用 base64格式，使用postman 的时候，记得把 + / 从 AES-online 换成 - _
gw55ZcBQOW+lgUmjzCRyzA==
gw55ZcBQOW-lgUmjzCRyzA==
  ┌────────────────────┬───────────────────┬────────────────────────────────────┐
  │        编码        │      字符集       │              对应 PHP              │
  ├────────────────────┼───────────────────┼────────────────────────────────────┤
  │ base64.StdEncoding │ A-Z a-z 0-9 + / = │ base64_encode()                    │
  ├────────────────────┼───────────────────┼────────────────────────────────────┤
  │ base64.URLEncoding │ A-Z a-z 0-9 - _ = │ strtr(base64_encode(), '+/', '-_') │
  └────────────────────┴───────────────────┴────────────────────────────────────┘
*/

type AbstractHelper struct {
}

func NewAbstractHelper() *AbstractHelper {
	return &AbstractHelper{}
}
