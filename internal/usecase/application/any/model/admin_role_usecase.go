package usecaseApplicationAnyModel

import (
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
)

type AdminRoleUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.AdminRoleModel
}

func NewAdminRoleUsecase(oAdminRoleModel outputPortAnyModel.AdminRoleModel, oAbstractUsecase *AbstractUsecase) usecasePortAnyModel.AdminRoleUsecase {
	return &AdminRoleUsecase{
		AbstractUsecase: oAbstractUsecase,
		AdminRoleModel:  oAdminRoleModel,
	}
}
