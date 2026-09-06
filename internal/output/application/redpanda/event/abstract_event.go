package outputApplicationRedpandaEvent

import (
	outputApplicationRedpanda "example/internal/output/application/redpanda"
)

// AbstractRedpanda 由 container wire 注入，Context / Client 等提升上去；
// event 這層目前不需要額外欄位，只是把 redpanda 的共用資源包一層。
type AbstractEvent struct {
	*outputApplicationRedpanda.AbstractRedpanda
}

func NewAbstractEvent(oAbstractRedpanda *outputApplicationRedpanda.AbstractRedpanda) *AbstractEvent {
	return &AbstractEvent{
		AbstractRedpanda: oAbstractRedpanda,
	}
}
