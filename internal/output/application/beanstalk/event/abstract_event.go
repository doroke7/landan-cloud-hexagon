package outputApplicationBeanstalkEvent

import (
	outputApplicationBeanstalk "example/internal/output/application/beanstalk"
)

// AbstractBeanstalk 由 container wire 注入，Context / Conn / Tube 提升上去；
// event 這層目前不需要額外欄位，只是把 beanstalk 的共用資源包一層。
type AbstractEvent struct {
	*outputApplicationBeanstalk.AbstractBeanstalk
}

func NewAbstractEvent(oAbstractBeanstalk *outputApplicationBeanstalk.AbstractBeanstalk) *AbstractEvent {
	return &AbstractEvent{
		AbstractBeanstalk: oAbstractBeanstalk,
	}
}
