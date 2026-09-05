package inputApplicationCronAdminResource

import (
	"go.uber.org/zap"

	pkgUtility "example/pkg/utility"

	inputApplicationCron "example/internal/input/application/cron"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
)

type AppUserHandler struct {
	*inputApplicationCron.AbstractHandler
	AdminResourceAppUserUsecase usecasePortAnyAdminResource.AppUserUsecase
}

func NewAppUserHandler(oAppUserUsecase usecasePortAnyAdminResource.AppUserUsecase, oAbstractHandler *inputApplicationCron.AbstractHandler) *AppUserHandler {
	return &AppUserHandler{
		AbstractHandler:             oAbstractHandler,
		AdminResourceAppUserUsecase: oAppUserUsecase,
	}
}

func (oSelf *AppUserHandler) IncreaseBalance() {
	err := oSelf.AdminResourceAppUserUsecase.IncreaseBalance(1, 10)
	if err != nil {
		pkgUtility.Logger(pkgUtility.Cron).Error("IncreaseBalance 失敗",
			zap.Error(err),
		)
		return
	}

	pkgUtility.Logger(pkgUtility.Cron).Info("IncreaseBalance 成功")
}
