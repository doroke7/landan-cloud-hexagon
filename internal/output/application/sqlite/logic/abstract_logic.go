package outputApplicationSqliteLogic

import (
	outputApplicationSqlite "example/internal/output/application/sqlite"
)

// AbstractSqlite 由 container wire 注入，DB / Context / Aop / FiltersToWheres 等提升上去；
// logic 這層目前不需要額外欄位，只是把 sqlite 的共用資源包一層。
type AbstractLogic struct {
	*outputApplicationSqlite.AbstractSqlite
}

func NewAbstractLogic(oAbstractSqlite *outputApplicationSqlite.AbstractSqlite) *AbstractLogic {
	return &AbstractLogic{
		AbstractSqlite: oAbstractSqlite,
	}
}
