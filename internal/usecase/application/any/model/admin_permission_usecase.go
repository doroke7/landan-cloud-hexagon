package usecaseApplicationAnyModel

import (
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
)

type AdminPermissionUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.AdminPermissionModel
}

func NewAdminPermissionUsecase(oAdminPermissionModel outputPortAnyModel.AdminPermissionModel, oAbstractUsecase *AbstractUsecase) usecasePortAnyModel.AdminPermissionUsecase {
	return &AdminPermissionUsecase{
		AbstractUsecase:      oAbstractUsecase,
		AdminPermissionModel: oAdminPermissionModel,
	}
}
