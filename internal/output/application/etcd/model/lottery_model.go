package outputApplicationEtcdModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
)

// LotteryModel 直接讀寫 etcd，不包其他 repository。
type LotteryModel struct {
	*AbstractModel
}

func NewLotteryModel(oAbstractModel *AbstractModel) outputPortAnyModel.LotteryModel {
	return &LotteryModel{
		AbstractModel: oAbstractModel,
	}
}

// WatchOneByKey 是讀：直接讀 etcd，沒有就回傳錯誤（不再自己生資料）。
func (oSelf *LotteryModel) WatchOneByKey(sKey string) (*domain.Lottery, error) {
	var oLottery domain.Lottery
	sCacheKey := oSelf.cacheKey(sKey)
	if oErr := oSelf.EtcdHelper.ReadCache(sCacheKey, &oLottery); oErr != nil {
		return nil, oErr
	}

	return &oLottery, nil
}

// EditOneByKey 是寫：把呼叫端給的 oValue 寫進 etcd。
func (oSelf *LotteryModel) EditOneByKey(oValue *domain.LotteryVariable, sKey string) error {

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
	if oErr := oSelf.EtcdHelper.WriteCache(sCacheKey, &oLottery); oErr != nil {
		return oErr
	}

	return nil
}

func (oSelf *LotteryModel) cacheKey(sKey string) string {
	return "lottery:" + sKey
}
