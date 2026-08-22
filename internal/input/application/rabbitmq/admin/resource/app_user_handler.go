package rabbitmq

import (
	"encoding/json"

	"go.uber.org/zap"

	amqp "github.com/rabbitmq/amqp091-go"

	pkg "example/pkg"

	inputApplicationRabbitmq "example/internal/input/application/rabbitmq"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
)

type increaseBalanceMessage struct {
	Id     uint `json:"id"`
	Amount uint `json:"amount"`
}

type AppUserHandler struct {
	*inputApplicationRabbitmq.AbstractHandler
	appUserUsecase usecasePortAnyAdminResource.AppUserUsecase
}

func NewAppUserHandler(oAppUserUsecase usecasePortAnyAdminResource.AppUserUsecase, oAbstractHandler *inputApplicationRabbitmq.AbstractHandler) *AppUserHandler {
	return &AppUserHandler{
		AbstractHandler: oAbstractHandler,
		appUserUsecase:  oAppUserUsecase,
	}
}

func (oSelf *AppUserHandler) IncreaseBalance(msg amqp.Delivery) {
	var payload increaseBalanceMessage
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		pkg.Logger(pkg.Consumer).Error("IncreaseBalance 訊息格式錯誤",
			zap.Error(err),
		)
		msg.Nack(false, false)
		return
	}

	bResult, err := oSelf.appUserUsecase.IncreaseBalance(payload.Id, payload.Amount)
	if err != nil {
		pkg.Logger(pkg.Consumer).Error("IncreaseBalance 失敗",
			zap.Uint("id", payload.Id),
			zap.Error(err),
		)
		msg.Nack(false, true)
		return
	}

	pkg.Logger(pkg.Consumer).Info("IncreaseBalance 成功",
		zap.Uint("id", payload.Id),
		zap.Bool("status", bResult),
	)
	msg.Ack(false)
}
