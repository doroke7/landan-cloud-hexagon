package outputApplicationMemoryLogic

import (
	outputApplicationMemory "example/internal/output/application/memory"
)

// AbstractMemory 由 container wire 注入，Context 提升上去；
// logic 這層目前不需要額外欄位，只是把 memory 的共用資源包一層。
type AbstractLogic struct {
	*outputApplicationMemory.AbstractMemory
}

func NewAbstractLogic(oAbstractMemory *outputApplicationMemory.AbstractMemory) *AbstractLogic {
	return &AbstractLogic{
		AbstractMemory: oAbstractMemory,
	}
}
