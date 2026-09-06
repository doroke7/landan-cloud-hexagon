package outputApplicationOracleLogic

import (
	outputApplicationOracle "example/internal/output/application/oracle"
)

// AbstractOracle 由 container wire 注入，DB / Context / Aop / FiltersToWheres 等提升上去；
// logic 這層目前不需要額外欄位，只是把 oracle 的共用資源包一層。
type AbstractLogic struct {
	*outputApplicationOracle.AbstractOracle
}

func NewAbstractLogic(oAbstractOracle *outputApplicationOracle.AbstractOracle) *AbstractLogic {
	return &AbstractLogic{
		AbstractOracle: oAbstractOracle,
	}
}
