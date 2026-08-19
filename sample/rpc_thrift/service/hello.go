package service

import (
	"context"
	"fmt"
)

func NewHello() *Hello {
	return &Hello{}
}

type Hello struct{}

func (oSelf *Hello) TestHello(ctx context.Context, sReq string) (string, error) {
	fmt.Printf("[Server] 收到客戶端請求: %s\n", sReq)
	return "Hello " + sReq + " (from Go Thrift Server)", nil
}
