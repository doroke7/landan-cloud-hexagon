package outputApplicationActivemqEvent

import (
	outputApplicationActivemq "example/internal/output/application/activemq"
)

// AbstractActivemq 由 container wire 注入，Context / Conn 方法（欄位）提升上去；
// event 這層目前不需要額外欄位，只是把 activemq 的共用資源包一層。
type AbstractEvent struct {
	*outputApplicationActivemq.AbstractActivemq
}

func NewAbstractEvent(oAbstractActivemq *outputApplicationActivemq.AbstractActivemq) *AbstractEvent {
	return &AbstractEvent{
		AbstractActivemq: oAbstractActivemq,
	}
}
