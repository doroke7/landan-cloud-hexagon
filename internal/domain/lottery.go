package domain

type Lottery struct {
	Id      uint   `json:"id"`
	Round   string `json:"round"`
	Time    int64  `json:"time"`
	Numbers string `json:"numbers"`
}

type LotteryRow struct {
	Id      uint   `json:"id"`
	Round   string `json:"round"`
	Time    int64  `json:"time"`
	Numbers string `json:"numbers"`
}

type LotteryValue struct {
	Round   *string `json:"round,omitempty"`
	Time    *int64  `json:"time,omitempty"`
	Numbers *string `json:"numbers,omitempty"`
}
