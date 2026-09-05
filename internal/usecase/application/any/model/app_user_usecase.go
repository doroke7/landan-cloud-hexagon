package usecaseApplicationAnyModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
)

type AppUserUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.AppUserModel
}

func NewAppUserUsecase(oAppUserRepository outputPortAnyModel.AppUserModel, oAbstractUsecase *AbstractUsecase) usecasePortAnyModel.AppUserUsecase {
	return &AppUserUsecase{
		AbstractUsecase: oAbstractUsecase,
		AppUserModel:    oAppUserRepository,
	}
}

func (oSelf *AppUserUsecase) ShowOneByName(sName string) (*domain.AppUser, error) {

	oAppUser, err := oSelf.AppUserModel.ShowOneByName(sName)

	return oAppUser, err
}

func (oSelf *AppUserUsecase) ShowOneById(iId uint) (*domain.AppUser, error) {

	oAppUser, err := oSelf.AppUserModel.ShowOneById(iId)

	return oAppUser, err
}

func (oSelf *AppUserUsecase) AddOne(oAppUser *domain.AppUserValue) (bool, error) {

	oErr := oSelf.AppUserModel.AddOne(oAppUser)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AppUserUsecase) IncreaseBalance(iId uint, iAmount uint) (bool, error) {

	oErr := oSelf.AppUserModel.IncreaseBalance(iId, iAmount)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}
