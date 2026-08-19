package resource

import (
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
)

type AppUserUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.AppUserModel
}

func NewAppUserUsecase(oAppUserModel outputPortAnyModel.AppUserModel, oAbstractUsecase *AbstractUsecase) usecasePortAnyAdminResource.AppUserUsecase {
	return &AppUserUsecase{
		AbstractUsecase: oAbstractUsecase,
		AppUserModel:    oAppUserModel,
	}
}

func (oSelf *AppUserUsecase) IncreaseBalance(id uint, amount uint) (bool, error) {

	_, oErr := oSelf.AppUserModel.IncreaseBalance(id, amount)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}
