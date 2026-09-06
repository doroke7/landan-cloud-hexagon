package outputApplicationMongodbModel

import (
	outputApplicationMongodb "example/internal/output/application/mongodb"
)

// AbstractMongodb 由 container wire 注入，Context / Client / Database / FiltersToFilter 等提升上去；
// model 這層目前不需要額外欄位，只是把 mongodb 的共用資源包一層。
type AbstractModel struct {
	*outputApplicationMongodb.AbstractMongodb
}

func NewAbstractModel(oAbstractMongodb *outputApplicationMongodb.AbstractMongodb) *AbstractModel {
	return &AbstractModel{
		AbstractMongodb: oAbstractMongodb,
	}
}
