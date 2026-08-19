package announcement

import (
	"log"
	"strconv"
	"strings"
	"time"

	pbSourceAnnouncement "example/pb/source/announcement"

	"google.golang.org/grpc"

	inputApplicationSource "example/internal/input/application/source"

	usecasePortAnyAnnoucement "example/internal/usecase/port/any/announcement"
)

type LotteryHandler struct {
	pbSourceAnnouncement.UnimplementedLotteryControllerServer
	*inputApplicationSource.AbstractHandler
	usecasePortAnyAnnoucement.LotteryUsecase
}

func NewLotteryHandler(oAbstractHandler *inputApplicationSource.AbstractHandler, oLotteryUsecase usecasePortAnyAnnoucement.LotteryUsecase) *LotteryHandler {
	return &LotteryHandler{
		AbstractHandler: oAbstractHandler,
		LotteryUsecase:  oLotteryUsecase,
	}
}

func (oSelf *LotteryHandler) Watch(oReq *pbSourceAnnouncement.LotteryWatchRequest, oStream grpc.ServerStreamingServer[pbSourceAnnouncement.LotteryWatchResponse]) error {

	oTicker := time.NewTicker(3 * time.Second)
	defer oTicker.Stop()

	var iCount int32 = 1

	for {
		select {
		// 情境 A：Client 斷線或 Cancel 請求
		case <-oStream.Context().Done():
			return oStream.Context().Err()

		// 情境 B：時間到，準備發送資料
		case <-oTicker.C:
			iCount++

			oLottert, oError := oSelf.LotteryUsecase.WatchOneByKey("SGS")

			if oError != nil {
				log.Printf("取得資料失敗: %v", oError)

			}

			var aNumbers []int32
			for _, sN := range strings.Split(oLottert.Numbers, ",") {
				iN, _ := strconv.Atoi(sN)
				aNumbers = append(aNumbers, int32(iN))
			}

			// 模擬產生開獎資料 (實務上這裡會是呼叫 UseCase / DB 撈資料)
			oResponse := &pbSourceAnnouncement.LotteryWatchResponse{
				Id:      int32(oLottert.Id),
				Round:   oLottert.Round,
				Time:    oLottert.Time,
				Numbers: aNumbers,
			}

			// 透過 stream.Send() 推送資料給 Client
			if oError := oStream.Send(oResponse); oError != nil {
				log.Printf("資料推送失敗: %v", oError)
				return oError // 發送失敗（通常是網路斷開），中斷 return
			}

			log.Printf("成功推送期號: %s", oResponse.Round)
		}
	}
}
