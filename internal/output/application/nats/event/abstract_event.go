package outputApplicationNatsEvent

import (
	outputApplicationNats "example/internal/output/application/nats"
)

// AbstractNats 由 container wire 注入，Context / Conn 提升上去；
// event 這層目前不需要額外欄位，只是把 nats 的共用資源包一層。
type AbstractEvent struct {
	*outputApplicationNats.AbstractNats
}

func NewAbstractEvent(oAbstractNats *outputApplicationNats.AbstractNats) *AbstractEvent {
	return &AbstractEvent{
		AbstractNats: oAbstractNats,
	}
}
