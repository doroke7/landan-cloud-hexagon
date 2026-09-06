package outputApplicationResourceEvent

import (
	"context"

	client "example/internal/client"
	outputApplicationResource "example/internal/output/application/resource"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractRepository 做法一致。
// AbstractResource 由 container wire 注入，ToFilters / ToSorters / ToPagination 方法提升上去；
// client 欄位保留在這裡。
type AbstractEvent struct {
	*outputApplicationResource.AbstractResource
	Context             context.Context
	ResourceLogicClient *client.Logic
	ResourceModelClient *client.Model
	ResourceEventClient *client.Event
}

func NewAbstractEvent(oContext context.Context, oResourceClient *client.ResourceClient, oAbstractResource *outputApplicationResource.AbstractResource) *AbstractEvent {
	return &AbstractEvent{
		AbstractResource:    oAbstractResource,
		Context:             oContext,
		ResourceLogicClient: oResourceClient.Logic,
		ResourceModelClient: oResourceClient.Model,
		ResourceEventClient: oResourceClient.Event,
	}
}
