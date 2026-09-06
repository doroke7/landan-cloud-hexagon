package outputApplicationCacheLogic

import (
	outputApplicationCache "example/internal/output/application/cache"
)

// AbstractCache 由 container wire 注入，CacheHelper / Context 提升上去；
// logic 這層目前不需要額外欄位，只是把 cache 的共用資源包一層。
type AbstractLogic struct {
	*outputApplicationCache.AbstractCache
}

func NewAbstractLogic(oAbstractCache *outputApplicationCache.AbstractCache) *AbstractLogic {
	return &AbstractLogic{
		AbstractCache: oAbstractCache,
	}
}
