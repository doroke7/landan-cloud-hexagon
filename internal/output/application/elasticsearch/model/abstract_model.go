package outputApplicationElasticsearchModel

import (
	outputApplicationElasticsearch "example/internal/output/application/elasticsearch"
)

// AbstractElasticsearch 由 container wire 注入，Client / Context / IndexName 等提升上去；
// model 這層目前不需要額外欄位，只是把 elasticsearch 的共用資源包一層。
type AbstractModel struct {
	*outputApplicationElasticsearch.AbstractElasticsearch
}

func NewAbstractModel(oAbstractElasticsearch *outputApplicationElasticsearch.AbstractElasticsearch) *AbstractModel {
	return &AbstractModel{
		AbstractElasticsearch: oAbstractElasticsearch,
	}
}
