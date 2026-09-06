package domain

/*
1. 綜合來說，不建議 proto = domain， 不好改動
2. 建議 保留 gRPC oProtoAdminUser 對 oDomainAdminUser 的轉換，
3. 如果為了節省轉換性能 就 proto domain 合一，增加的性能只有一點，卻增加很大的改動難度
      （譬如， 如果對方的接口格式改變了，我們需要跟著改，那內部domain 也得改）

*/

type Lottery struct {
	Id      uint64 `json:"id"`
	Round   string `json:"round"`
	Time    uint64 `json:"time"`
	Numbers string `json:"numbers"`
}

type LotteryValue struct {
	Round   *string `json:"round,omitempty"`
	Time    *uint64 `json:"time,omitempty"`
	Numbers *string `json:"numbers,omitempty"`
}
