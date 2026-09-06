package outputApplicationGaussdbModel

import (
	outputApplicationGaussdb "example/internal/output/application/gaussdb"
)

// AbstractGaussdb 由 container wire 注入，DB / Context / Aop / FiltersToWheres 等提升上去；
// model 這層目前不需要額外欄位，只是把 gaussdb 的共用資源包一層。
type AbstractModel struct {
	*outputApplicationGaussdb.AbstractGaussdb
}

func NewAbstractModel(oAbstractGaussdb *outputApplicationGaussdb.AbstractGaussdb) *AbstractModel {
	return &AbstractModel{
		AbstractGaussdb: oAbstractGaussdb,
	}
}
