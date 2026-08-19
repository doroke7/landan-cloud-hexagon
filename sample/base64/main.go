package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
)

func Base64ToURLBase64(s string) string {
	// 標準 Base64 decode
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}

	// URL Base64 encode（不帶 padding）
	return base64.RawURLEncoding.EncodeToString(data)
}

func main() {
	origin := "SGVsbG8rLw==" // Hello+/

	result1 := Base64ToURLBase64(origin)
	result2 := url.QueryEscape(origin)

	fmt.Println(result1)
	fmt.Println(result2)

}
