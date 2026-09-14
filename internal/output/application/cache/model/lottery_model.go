package outputApplicationCacheModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
)

// LotteryModel 直接讀寫 redis，不包其他 repository。
type LotteryModel struct {
	*AbstractModel
}

func NewLotteryModel(oAbstractModel *AbstractModel) outputPortAnyModel.LotteryModel {
	return &LotteryModel{
		AbstractModel: oAbstractModel,
	}
}

// WatchOneByKey 是讀：直接讀 redis，沒有就回傳錯誤（不再自己生資料）。
func (oSelf *LotteryModel) WatchOneByKey(sKey string) (*domain.Lottery, error) {
	var oLottery domain.Lottery
	sCacheKey := oSelf.cacheKey(sKey)
	if err := oSelf.CacheHelper.ReadCache(sCacheKey, &oLottery); err != nil {
		return nil, err
	}

	return &oLottery, nil
}

// EditOneByKey 是寫：把呼叫端給的 oValue 寫進 redis。
func (oSelf *LotteryModel) EditOneByKey(oValue *domain.LotteryValue, sKey string) error {

	var oLottery domain.Lottery

	if oValue.Round != nil {
		oLottery.Round = *oValue.Round
	}

	if oValue.Time != nil {
		oLottery.Time = *oValue.Time
	}

	if oValue.Numbers != nil {
		oLottery.Numbers = *oValue.Numbers
	}

	sCacheKey := oSelf.cacheKey(sKey)
	if err := oSelf.CacheHelper.WriteCache(sCacheKey, &oLottery); err != nil {
		return err
	}

	return nil
}

func (oSelf *LotteryModel) cacheKey(sKey string) string {
	return "lottery:" + sKey
}
