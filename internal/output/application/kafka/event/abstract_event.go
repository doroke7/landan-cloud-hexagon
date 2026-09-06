package outputApplicationKafkaEvent

import (
	outputApplicationKafka "example/internal/output/application/kafka"
)

// AbstractKafka 由 container wire 注入，Context / Client / Producer 提升上去；
// event 這層目前不需要額外欄位，只是把 kafka 的共用資源包一層。
type AbstractEvent struct {
	*outputApplicationKafka.AbstractKafka
}

func NewAbstractEvent(oAbstractKafka *outputApplicationKafka.AbstractKafka) *AbstractEvent {
	return &AbstractEvent{
		AbstractKafka: oAbstractKafka,
	}
}
