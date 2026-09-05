package usecaseApplicationAnyModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
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

func (oSelf *AdminPermissionUsecase) AddOne(oAdminPermission *domain.AdminPermissionValue) error {

	oErr := oSelf.AdminPermissionModel.AddOne(oAdminPermission)

	return oErr
}

func (oSelf *AdminPermissionUsecase) ShowOneById(iId uint64) (*domain.AdminPermission, error) {

	oAdminPermission, oErr := oSelf.AdminPermissionModel.ShowOneById(iId)

	return oAdminPermission, oErr
}

func (oSelf *AdminPermissionUsecase) EditOneById(oAdminPermission *domain.AdminPermissionValue, iId uint64) error {

	oErr := oSelf.AdminPermissionModel.EditOneById(oAdminPermission, iId)

	return oErr
}

func (oSelf *AdminPermissionUsecase) RemoveOneById(iId uint64) error {

	oErr := oSelf.AdminPermissionModel.RemoveOneById(iId)

	return oErr
}

func (oSelf *AdminPermissionUsecase) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	iTotal, oErr := oSelf.AdminPermissionModel.TotalByFilters(aFilters)

	return uint64(iTotal), oErr
}
