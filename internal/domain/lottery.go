package domain

import bootstrap "example/bootstrap"

type Lottery struct {
	Id      uint   `json:"id"`
	Round   string `json:"round"`
	Time    int64  `json:"time"`
	Numbers string `json:"numbers"`
}

func (LotteryRow) TableName() string {
	return bootstrap.CONFIG.DATABASE.PREFIX + "lotteries"
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
