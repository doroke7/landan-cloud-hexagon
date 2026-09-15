package usecaseApplicationAnyLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
)

type AdminRoleUsecase struct {
	*AbstractUsecase
	outputPortAnyLogic.AdminRoleLogic
}

func NewAdminRoleUsecase(oAbstractUsecase *AbstractUsecase, oAdminRoleLogic outputPortAnyLogic.AdminRoleLogic) usecasePortAnyLogic.AdminRoleUsecase {
	return &AdminRoleUsecase{
		AbstractUsecase: oAbstractUsecase,
		AdminRoleLogic:  oAdminRoleLogic,
	}
}

func (oSelf *AdminRoleUsecase) AddAdminRole(oVariable *domain.AdminRoleVariable) error {

	oErr := oSelf.AdminRoleLogic.AddAdminRole(oVariable)

	return oErr
}
