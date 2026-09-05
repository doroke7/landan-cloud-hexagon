package usecaseApplicationAnyAdminResource

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkgInput "example/pkg/input"
)

type AdminPermissionUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	outputPortAnyModel.AdminPermissionModel
}

func NewAdminPermissionUsecase(oAdminPermissionModel outputPortAnyModel.AdminPermissionModel, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminResource.AdminPermissionUsecase {
	return &AdminPermissionUsecase{
		AbstractUsecase:      oAbstractUsecase,
		AdminPermissionModel: oAdminPermissionModel,
	}
}

func (oSelf *AdminPermissionUsecase) AddOne(oAdminPermission *domain.AdminPermissionValue) error {

	oErr := oSelf.AdminPermissionModel.AddOne(oAdminPermission)

	return oErr
}

func (oSelf *AdminPermissionUsecase) ShowOne(iId uint) (*domain.AdminPermission, error) {

	oAdminPermission, oErr := oSelf.AdminPermissionModel.ShowOneById(iId)

	return oAdminPermission, oErr
}

func (oSelf *AdminPermissionUsecase) EditOne(oAdminPermission *domain.AdminPermissionValue, iId uint) (bool, error) {

	oErr := oSelf.AdminPermissionModel.EditOneById(oAdminPermission, iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AdminPermissionUsecase) RemoveOne(iId uint) (bool, error) {

	oErr := oSelf.AdminPermissionModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AdminPermissionUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermission, uint, error) {

	aAdminPermissions, oErr := oSelf.AdminPermissionModel.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, 0, oErr
	}

	iTotal, oErr := oSelf.AdminPermissionModel.TotalByFilters(aFilters)

	return aAdminPermissions, iTotal, oErr
}
