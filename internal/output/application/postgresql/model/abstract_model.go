package outputApplicationPostgresqlModel

import (
	outputApplicationPostgresql "example/internal/output/application/postgresql"
)

// AbstractPostgresql 由 container wire 注入，DB / Context / Aop / FiltersToWheres 等提升上去；
// model 這層目前不需要額外欄位，只是把 postgresql 的共用資源包一層。
type AbstractModel struct {
	*outputApplicationPostgresql.AbstractPostgresql
}

func NewAbstractModel(oAbstractPostgresql *outputApplicationPostgresql.AbstractPostgresql) *AbstractModel {
	return &AbstractModel{
		AbstractPostgresql: oAbstractPostgresql,
	}
}
