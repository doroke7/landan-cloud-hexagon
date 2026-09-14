package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
)

type AppUserLogic struct {
	*AbstractLogic
}

func NewAppUserLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AppUserLogic {
	oLogic := &AppUserLogic{
		AbstractLogic: oAbstractLogic,
	}

	return oLogic
}

func (oSelf *AppUserLogic) AddAppUser(oAppUserValue *domain.AppUserValue) (*domain.AppUser, error) {

	oAppUser := &domain.AppUser{
		Id:       1,
		Name:     "11",
		Password: "222222",
	}

	return oAppUser, nil
}
