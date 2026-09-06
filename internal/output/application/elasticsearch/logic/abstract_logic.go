package outputApplicationElasticsearchLogic

import (
	outputApplicationElasticsearch "example/internal/output/application/elasticsearch"
)

// AbstractElasticsearch 由 container wire 注入，Client / Context / IndexName 等提升上去；
// logic 這層目前不需要額外欄位，只是把 elasticsearch 的共用資源包一層。
type AbstractLogic struct {
	*outputApplicationElasticsearch.AbstractElasticsearch
}

func NewAbstractLogic(oAbstractElasticsearch *outputApplicationElasticsearch.AbstractElasticsearch) *AbstractLogic {
	return &AbstractLogic{
		AbstractElasticsearch: oAbstractElasticsearch,
	}
}
