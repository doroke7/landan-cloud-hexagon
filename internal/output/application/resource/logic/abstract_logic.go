package resource

import (
	"context"

	client "example/internal/client"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractRepository 做法一致。
type AbstractLogic struct {
	Context             context.Context
	ResourceLogicClient *client.Logic
}

func NewAbstractLogic(oContext context.Context, oResourceClient *client.ResourceClient) *AbstractLogic {
	return &AbstractLogic{
		Context:             oContext,
		ResourceLogicClient: oResourceClient.Logic,
	}
}
