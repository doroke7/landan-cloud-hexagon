package outputApplicationZeromqEvent

import (
	outputApplicationZeromq "example/internal/output/application/zeromq"
)

// AbstractZeromq 由 container wire 注入，Context / Socket 提升上去；
// event 這層目前不需要額外欄位，只是把 zeromq 的共用資源包一層。
type AbstractEvent struct {
	*outputApplicationZeromq.AbstractZeromq
}

func NewAbstractEvent(oAbstractZeromq *outputApplicationZeromq.AbstractZeromq) *AbstractEvent {
	return &AbstractEvent{
		AbstractZeromq: oAbstractZeromq,
	}
}
