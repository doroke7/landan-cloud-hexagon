package outputApplicationResourceEvent

import (
	"context"

	client "example/internal/client"
	outputApplicationResource "example/internal/output/application/resource"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractRepository 做法一致。
// AbstractResource 由 container wire 注入，ToFilters / ToSorters / ToPagination 方法提升上去；
// event 這層只會用到 ResourceEventClient。
type AbstractEvent struct {
	*outputApplicationResource.AbstractResource
	Context             context.Context
	ResourceEventClient *client.Event
}

func NewAbstractEvent(oContext context.Context, oResourceClient *client.ResourceClient, oAbstractResource *outputApplicationResource.AbstractResource) *AbstractEvent {
	oEvent := &AbstractEvent{
		AbstractResource:    oAbstractResource,
		Context:             oContext,
		ResourceEventClient: oResourceClient.Event,
	}

	return oEvent
}
