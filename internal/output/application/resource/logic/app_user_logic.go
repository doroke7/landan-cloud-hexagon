package resource

import (
	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AppUserLogic struct {
	*resourceBase.AbstractResource
}

func NewAppUserLogic(oAbstractLogic *resourceBase.AbstractResource) outputPortAnyModel.AppUserModel {
	return &AppUserLogic{
		AbstractResource: oAbstractLogic,
	}
}

func (oSelf *AppUserLogic) AddAppUser(oAppUser *domain.AppUser) (*domain.AppUser, error) {

	return &domain.AppUser{
		Id:       1,
		Name:     "11",
		Password: "222222",
	}, nil
}

func (oSelf *AppUserLogic) IncreaseBalance(id uint, amount uint) (bool, error) {

	return true, nil
}
