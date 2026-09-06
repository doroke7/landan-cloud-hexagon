package outputApplicationTidbLogic

import (
	outputApplicationTidb "example/internal/output/application/tidb"
)

// AbstractTidb 由 container wire 注入，DB / Context / Aop / FiltersToWheres 等提升上去；
// logic 這層目前不需要額外欄位，只是把 tidb 的共用資源包一層。
type AbstractLogic struct {
	*outputApplicationTidb.AbstractTidb
}

func NewAbstractLogic(oAbstractTidb *outputApplicationTidb.AbstractTidb) *AbstractLogic {
	return &AbstractLogic{
		AbstractTidb: oAbstractTidb,
	}
}
