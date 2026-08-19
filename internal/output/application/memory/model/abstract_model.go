package model

import (
	"context"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache.AbstractModel 的做法一致。
type AbstractModel struct {
	Context context.Context
}

func NewAbstractModel(oContext context.Context) *AbstractModel {
	return &AbstractModel{
		Context: oContext,
	}
}
