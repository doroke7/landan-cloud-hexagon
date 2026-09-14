package outputApplicationMemoryModel

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type LotteryModel struct {
	*AbstractModel
	startTime time.Time
}

func NewLotteryModel(oAbstractModel *AbstractModel) outputPortAnyModel.LotteryModel {
	return &LotteryModel{
		AbstractModel: oAbstractModel,
		startTime:     time.Now(),
	}
}

// WatchOneByKey 是讀：依照從 startTime 到現在經過的分鐘數計算目前這一期。
func (oSelf *LotteryModel) WatchOneByKey(sKey string) (*domain.Lottery, error) {

	oDuration := time.Since(oSelf.startTime)
	fMinutes := oDuration.Minutes()
	iId := uint(fMinutes) + 1

	iCount := 4 // 產生幾個數字

	aNumbers := make([]string, iCount)

	for i := 0; i < iCount; i++ {
		n := rand.IntN(99) + 1 // 1~99
		aNumbers[i] = strconv.Itoa(n)
	}

	sNumbers := strings.Join(aNumbers, ",")

	sRound := fmt.Sprintf("2026-%04d", iId)
	oNow := time.Now()
	iTime := uint64(oNow.UnixNano())

	oLottery := &domain.Lottery{
		Id:      uint64(iId),
		Round:   sRound,
		Time:    iTime,
		Numbers: sNumbers,
	}

	return oLottery, nil
}

// EditOneByKey 是寫，但 memory 沒有真的儲存空間可以落地，
// 這裡只是回傳成功，不做任何持久化（demo 用途）。
func (oSelf *LotteryModel) EditOneByKey(oValue *domain.LotteryVariable, sKey string) error {
	return nil
}
