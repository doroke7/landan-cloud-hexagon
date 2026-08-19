package resource

import (
	"context"

	client "example/internal/client"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractRepository 做法一致。
type AbstractModel struct {
	Context             context.Context
	ResourceModelClient *client.Model
}

func NewAbstractModel(oContext context.Context, oResourceClient *client.ResourceClient) *AbstractModel {
	return &AbstractModel{
		Context:             oContext,
		ResourceModelClient: oResourceClient.Model,
	}
}
