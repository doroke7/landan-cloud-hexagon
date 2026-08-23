package outputApplicationEtcdModel

import (
	domain "example/internal/domain"
	etcdBase "example/internal/output/application/etcd"
	outputPortAnyModel "example/internal/output/port/any/model"
)

// LotteryModel 直接讀寫 etcd，不包其他 repository。
type LotteryModel struct {
	*etcdBase.AbstractEtcd
}

func NewLotteryModel(oAbstractEtcd *etcdBase.AbstractEtcd) outputPortAnyModel.LotteryModel {
	return &LotteryModel{
		AbstractEtcd: oAbstractEtcd,
	}
}

// WatchOneByKey 是讀：直接讀 etcd，沒有就回傳錯誤（不再自己生資料）。
func (oSelf *LotteryModel) WatchOneByKey(sKey string) (*domain.Lottery, error) {
	var oLottery domain.Lottery
	if oErr := oSelf.EtcdHelper.ReadCache(oSelf.cacheKey(sKey), &oLottery); oErr != nil {
		return nil, oErr
	}

	return &oLottery, nil
}

// EditOneByKey 是寫：把呼叫端給的 oValue 寫進 etcd。
func (oSelf *LotteryModel) EditOneByKey(oValue *domain.LotteryValue, sKey string) (bool, error) {

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

	if oErr := oSelf.EtcdHelper.WriteCache(oSelf.cacheKey(sKey), &oLottery); oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *LotteryModel) cacheKey(sKey string) string {
	return "lottery:" + sKey
}
