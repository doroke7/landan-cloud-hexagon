package outputApplicationSqliteModel

import (
	outputApplicationSqlite "example/internal/output/application/sqlite"
)

// AbstractSqlite 由 container wire 注入，DB / Context / Aop / FiltersToWheres 等提升上去；
// model 這層目前不需要額外欄位，只是把 sqlite 的共用資源包一層。
type AbstractModel struct {
	*outputApplicationSqlite.AbstractSqlite
}

func NewAbstractModel(oAbstractSqlite *outputApplicationSqlite.AbstractSqlite) *AbstractModel {
	return &AbstractModel{
		AbstractSqlite: oAbstractSqlite,
	}
}
