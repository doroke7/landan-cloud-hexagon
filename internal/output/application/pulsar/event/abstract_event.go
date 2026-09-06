package outputApplicationPulsarEvent

import (
	outputApplicationPulsar "example/internal/output/application/pulsar"
)

// AbstractPulsar 由 container wire 注入，Context / Client 等提升上去；
// event 這層目前不需要額外欄位，只是把 pulsar 的共用資源包一層。
type AbstractEvent struct {
	*outputApplicationPulsar.AbstractPulsar
}

func NewAbstractEvent(oAbstractPulsar *outputApplicationPulsar.AbstractPulsar) *AbstractEvent {
	return &AbstractEvent{
		AbstractPulsar: oAbstractPulsar,
	}
}
