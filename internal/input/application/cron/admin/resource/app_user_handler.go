package cron

import (
	"go.uber.org/zap"

	pkg "example/pkg"

	inputApplicationCron "example/internal/input/application/cron"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
)

type AppUserHandler struct {
	*inputApplicationCron.AbstractHandler
	appUserUsecase usecasePortAnyAdminResource.AppUserUsecase
}

func NewAppUserHandler(oAppUserUsecase usecasePortAnyAdminResource.AppUserUsecase, oAbstractHandler *inputApplicationCron.AbstractHandler) *AppUserHandler {
	return &AppUserHandler{
		AbstractHandler: oAbstractHandler,
		appUserUsecase:  oAppUserUsecase,
	}
}

func (oSelf *AppUserHandler) IncreaseBalance() {
	bResult, err := oSelf.appUserUsecase.IncreaseBalance(1, 10)
	if err != nil {
		pkg.Logger(pkg.Cron).Error("IncreaseBalance 失敗",
			zap.Error(err),
		)
		return
	}

	pkg.Logger(pkg.Cron).Info("IncreaseBalance 成功",
		zap.Bool("status", bResult),
	)
}
