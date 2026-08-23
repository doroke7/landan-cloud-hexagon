package outputApplicationMemory

import (
	"context"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache.AbstractCache 的做法一致。
type AbstractMemory struct {
	Context context.Context
}

func NewAbstractMemory(oContext context.Context) *AbstractMemory {
	return &AbstractMemory{
		Context: oContext,
	}
}
