package outputApplicationRabbitmqEvent

import (
	outputApplicationRabbitmq "example/internal/output/application/rabbitmq"
)

// AbstractRabbitmq 由 container wire 注入，Context / Conn / Channel 提升上去；
// event 這層目前不需要額外欄位，只是把 rabbitmq 的共用資源包一層。
type AbstractEvent struct {
	*outputApplicationRabbitmq.AbstractRabbitmq
}

func NewAbstractEvent(oAbstractRabbitmq *outputApplicationRabbitmq.AbstractRabbitmq) *AbstractEvent {
	return &AbstractEvent{
		AbstractRabbitmq: oAbstractRabbitmq,
	}
}
