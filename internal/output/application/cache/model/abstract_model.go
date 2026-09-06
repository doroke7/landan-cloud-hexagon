package outputApplicationCacheModel

import (
	outputApplicationCache "example/internal/output/application/cache"
)

// AbstractCache 由 container wire 注入，CacheHelper / Context 提升上去；
// model 這層目前不需要額外欄位，只是把 cache 的共用資源包一層。
type AbstractModel struct {
	*outputApplicationCache.AbstractCache
}

func NewAbstractModel(oAbstractCache *outputApplicationCache.AbstractCache) *AbstractModel {
	return &AbstractModel{
		AbstractCache: oAbstractCache,
	}
}
