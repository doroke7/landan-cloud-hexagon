package usecaseApplicationAnyLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
)

type AppUserUsecase struct {
	*AbstractUsecase
	outputPortAnyLogic.AppUserLogic
}

func NewAppUserUsecase(oAppUserRepository outputPortAnyLogic.AppUserLogic, oAbstractUsecase *AbstractUsecase) usecasePortAnyLogic.AppUserUsecase {
	return &AppUserUsecase{
		AbstractUsecase: oAbstractUsecase,
		AppUserLogic:    oAppUserRepository,
	}
}

func (oSelf *AppUserUsecase) AddAppUser(oAppUser *domain.AppUser) (bool, error) {

	return true, nil
}
