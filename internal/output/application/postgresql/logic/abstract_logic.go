package outputApplicationPostgresqlLogic

import (
	outputApplicationPostgresql "example/internal/output/application/postgresql"
)

// AbstractPostgresql 由 container wire 注入，DB / Context / Aop / FiltersToWheres 等提升上去；
// logic 這層目前不需要額外欄位，只是把 postgresql 的共用資源包一層。
type AbstractLogic struct {
	*outputApplicationPostgresql.AbstractPostgresql
}

func NewAbstractLogic(oAbstractPostgresql *outputApplicationPostgresql.AbstractPostgresql) *AbstractLogic {
	return &AbstractLogic{
		AbstractPostgresql: oAbstractPostgresql,
	}
}
