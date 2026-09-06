package outputApplicationClickhouseModel

import (
	outputApplicationClickhouse "example/internal/output/application/clickhouse"
)

// AbstractClickhouse 由 container wire 注入，DB / Context / Aop / FiltersToWheres 等提升上去；
// model 這層目前不需要額外欄位，只是把 clickhouse 的共用資源包一層。
type AbstractModel struct {
	*outputApplicationClickhouse.AbstractClickhouse
}

func NewAbstractModel(oAbstractClickhouse *outputApplicationClickhouse.AbstractClickhouse) *AbstractModel {
	return &AbstractModel{
		AbstractClickhouse: oAbstractClickhouse,
	}
}
