package outputApplicationMysqlModel

import (
	outputApplicationMysql "example/internal/output/application/mysql"
)

// AbstractMysql 由 container wire 注入，DB / Context / Aop / FiltersToWheres 等提升上去；
// model 這層目前不需要額外欄位，只是把 mysql 的共用資源包一層。
type AbstractModel struct {
	*outputApplicationMysql.AbstractMysql
}

func NewAbstractModel(oAbstractMysql *outputApplicationMysql.AbstractMysql) *AbstractModel {
	return &AbstractModel{
		AbstractMysql: oAbstractMysql,
	}
}
