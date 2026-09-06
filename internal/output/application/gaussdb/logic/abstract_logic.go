package outputApplicationGaussdbLogic

import (
	outputApplicationGaussdb "example/internal/output/application/gaussdb"
)

// AbstractGaussdb 由 container wire 注入，DB / Context / Aop / FiltersToWheres 等提升上去；
// logic 這層目前不需要額外欄位，只是把 gaussdb 的共用資源包一層。
type AbstractLogic struct {
	*outputApplicationGaussdb.AbstractGaussdb
}

func NewAbstractLogic(oAbstractGaussdb *outputApplicationGaussdb.AbstractGaussdb) *AbstractLogic {
	return &AbstractLogic{
		AbstractGaussdb: oAbstractGaussdb,
	}
}
