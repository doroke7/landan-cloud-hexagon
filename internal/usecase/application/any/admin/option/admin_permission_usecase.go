package usecaseApplicationAnyAdminOption

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminOption "example/internal/usecase/port/any/admin/option"
	pkgInput "example/pkg/input"
)

type AdminPermissionUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	AdminPermissionModel outputPortAnyModel.AdminPermissionModel
}

func NewAdminPermissionUsecase(oAdminPermissionModel outputPortAnyModel.AdminPermissionModel, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminOption.AdminPermissionUsecase {
	return &AdminPermissionUsecase{
		AbstractUsecase:      oAbstractUsecase,
		AdminPermissionModel: oAdminPermissionModel,
	}
}

func (oSelf *AdminPermissionUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermission, uint64, error) {

	aAdminPermissions, oErr := oSelf.AdminPermissionModel.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, 0, oErr
	}

	iTotal, oErr := oSelf.AdminPermissionModel.TotalByFilters(aFilters)

	return aAdminPermissions, uint64(iTotal), oErr
}
