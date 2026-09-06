package outputApplicationTidbModel

import (
	outputApplicationTidb "example/internal/output/application/tidb"
)

// AbstractTidb 由 container wire 注入，DB / Context / Aop / FiltersToWheres 等提升上去；
// model 這層目前不需要額外欄位，只是把 tidb 的共用資源包一層。
type AbstractModel struct {
	*outputApplicationTidb.AbstractTidb
}

func NewAbstractModel(oAbstractTidb *outputApplicationTidb.AbstractTidb) *AbstractModel {
	return &AbstractModel{
		AbstractTidb: oAbstractTidb,
	}
}
