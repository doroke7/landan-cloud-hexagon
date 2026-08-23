package usecaseApplicationAnyAdminResource

import (
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
)

type AppUserUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	outputPortAnyModel.AppUserModel
}

func NewAppUserUsecase(oAppUserModel outputPortAnyModel.AppUserModel, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminResource.AppUserUsecase {
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
