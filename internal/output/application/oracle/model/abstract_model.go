package outputApplicationOracleModel

import (
	outputApplicationOracle "example/internal/output/application/oracle"
)

// AbstractOracle 由 container wire 注入，DB / Context / Aop / FiltersToWheres 等提升上去；
// model 這層目前不需要額外欄位，只是把 oracle 的共用資源包一層。
type AbstractModel struct {
	*outputApplicationOracle.AbstractOracle
}

func NewAbstractModel(oAbstractOracle *outputApplicationOracle.AbstractOracle) *AbstractModel {
	return &AbstractModel{
		AbstractOracle: oAbstractOracle,
	}
}
