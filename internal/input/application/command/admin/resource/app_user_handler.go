package command

import (
	inputApplicationCommand "example/internal/input/application/command"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
)

type AppUserHandler struct {
	*inputApplicationCommand.AbstractHandler
	appUserUsecase usecasePortAnyAdminResource.AppUserUsecase
}

func NewAppUserHandler(oAppUserUsecase usecasePortAnyAdminResource.AppUserUsecase, oAbstractHandler *inputApplicationCommand.AbstractHandler) *AppUserHandler {
	return &AppUserHandler{
		AbstractHandler: oAbstractHandler,
		appUserUsecase:  oAppUserUsecase,
	}
}

// IncreaseBalance 只負責轉呼叫 usecase，不知道自己是被 CLI 呼叫的——
// cobra.Command 的組裝、container 的建立時機，都交給 cmd/register 那層決定。
func (oSelf *AppUserHandler) IncreaseBalance(id uint, amount uint) error {
	_, err := oSelf.appUserUsecase.IncreaseBalance(id, amount)
	return err
}
