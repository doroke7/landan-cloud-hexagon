package model

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	domain "example/internal/domain"
	memoryBase "example/internal/output/application/memory"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type LotteryModel struct {
	*memoryBase.AbstractMemory
	startTime time.Time
}

func NewLotteryModel(oAbstractMemory *memoryBase.AbstractMemory) outputPortAnyModel.LotteryModel {
	return &LotteryModel{
		AbstractMemory: oAbstractMemory,
		startTime:      time.Now(),
	}
}

// WatchOneByKey 是讀：依照從 startTime 到現在經過的分鐘數計算目前這一期。
func (oSelf *LotteryModel) WatchOneByKey(sKey string) (*domain.Lottery, error) {

	iId := uint(time.Since(oSelf.startTime).Minutes()) + 1

	iCount := 4 // 產生幾個數字

	aNumbers := make([]string, iCount)

	for i := 0; i < iCount; i++ {
		n := rand.IntN(99) + 1 // 1~99
		aNumbers[i] = strconv.Itoa(n)
	}

	sNumbers := strings.Join(aNumbers, ",")

	return &domain.Lottery{
		Id:      iId,
		Round:   fmt.Sprintf("2026-%04d", iId),
		Time:    time.Now().UnixNano(),
		Numbers: sNumbers,
	}, nil

}

// EditOneByKey 是寫，但 memory 沒有真的儲存空間可以落地，
// 這裡只是回傳成功，不做任何持久化（demo 用途）。
func (oSelf *LotteryModel) EditOneByKey(oValue *domain.LotteryValue, sKey string) (bool, error) {
	return true, nil
}
