package outputApplicationMemoryModel

import (
	outputApplicationMemory "example/internal/output/application/memory"
)

// AbstractMemory 由 container wire 注入，Context 提升上去；
// model 這層目前不需要額外欄位，只是把 memory 的共用資源包一層。
type AbstractModel struct {
	*outputApplicationMemory.AbstractMemory
}

func NewAbstractModel(oAbstractMemory *outputApplicationMemory.AbstractMemory) *AbstractModel {
	return &AbstractModel{
		AbstractMemory: oAbstractMemory,
	}
}
