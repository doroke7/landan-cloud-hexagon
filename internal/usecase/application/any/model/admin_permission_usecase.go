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

func (oSelf *AdminPermissionUsecase) AddOne(oAdminPermission *domain.AdminPermissionValue) (bool, error) {

	oErr := oSelf.AdminPermissionModel.AddOne(oAdminPermission)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AdminPermissionUsecase) EditOneById(oAdminPermission *domain.AdminPermissionValue, iId uint) (bool, error) {

	oErr := oSelf.AdminPermissionModel.EditOneById(oAdminPermission, iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AdminPermissionUsecase) RemoveOneById(iId uint) (bool, error) {

	oErr := oSelf.AdminPermissionModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AdminPermissionUsecase) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	iTotal, oErr := oSelf.AdminPermissionModel.TotalByFilters(aFilters)

	return uint64(iTotal), oErr
}
