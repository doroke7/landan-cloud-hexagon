package outputApplicationClickhouseLogic

import (
	outputApplicationClickhouse "example/internal/output/application/clickhouse"
)

// AbstractClickhouse 由 container wire 注入，DB / Context / Aop / FiltersToWheres 等提升上去；
// logic 這層目前不需要額外欄位，只是把 clickhouse 的共用資源包一層。
type AbstractLogic struct {
	*outputApplicationClickhouse.AbstractClickhouse
}

func NewAbstractLogic(oAbstractClickhouse *outputApplicationClickhouse.AbstractClickhouse) *AbstractLogic {
	return &AbstractLogic{
		AbstractClickhouse: oAbstractClickhouse,
	}
}
