package outputApplicationMosquittoEvent

import (
	outputApplicationMosquitto "example/internal/output/application/mosquitto"
)

// AbstractMosquitto 由 container wire 注入，Context / Client 提升上去；
// event 這層目前不需要額外欄位，只是把 mosquitto 的共用資源包一層。
type AbstractEvent struct {
	*outputApplicationMosquitto.AbstractMosquitto
}

func NewAbstractEvent(oAbstractMosquitto *outputApplicationMosquitto.AbstractMosquitto) *AbstractEvent {
	return &AbstractEvent{
		AbstractMosquitto: oAbstractMosquitto,
	}
}
