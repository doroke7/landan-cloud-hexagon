package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
)

type AppUserLogic struct {
	*AbstractLogic
}

func NewAppUserLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AppUserLogic {
	return &AppUserLogic{
		AbstractLogic: oAbstractLogic,
	}
}

func (oSelf *AppUserLogic) AddAppUser(oAppUserValue *domain.AppUserValue) (*domain.AppUser, error) {

	return &domain.AppUser{
		Id:       1,
		Name:     "11",
		Password: "222222",
	}, nil
}
