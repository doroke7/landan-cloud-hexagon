package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyLogic "example/internal/output/port/any/logic"
)

type AppUserLogic struct {
	*resourceBase.AbstractResource
}

func NewAppUserLogic(oAbstractLogic *resourceBase.AbstractResource) outputPortAnyLogic.AppUserLogic {
	return &AppUserLogic{
		AbstractResource: oAbstractLogic,
	}
}

func (oSelf *AppUserLogic) AddAppUser(oAppUserValue *domain.AppUserValue) (*domain.AppUser, error) {

	return &domain.AppUser{
		Id:       1,
		Name:     "11",
		Password: "222222",
	}, nil
}
