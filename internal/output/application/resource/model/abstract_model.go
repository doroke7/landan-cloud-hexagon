package outputApplicationResourceModel

import (
	"context"

	client "example/internal/client"
	outputApplicationResource "example/internal/output/application/resource"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractRepository 做法一致。
// AbstractResource 由 container wire 注入，ToFilters / ToSorters / ToPagination 方法提升上去；
// model 這層只會用到 ResourceModelClient。
type AbstractModel struct {
	*outputApplicationResource.AbstractResource
	Context             context.Context
	ResourceModelClient *client.Model
}

func NewAbstractModel(oContext context.Context, oResourceClient *client.ResourceClient, oAbstractResource *outputApplicationResource.AbstractResource) *AbstractModel {
	return &AbstractModel{
		AbstractResource:    oAbstractResource,
		Context:             oContext,
		ResourceModelClient: oResourceClient.Model,
	}
}
