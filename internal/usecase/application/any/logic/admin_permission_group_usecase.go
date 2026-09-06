package usecaseApplicationAnyLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupUsecase struct {
	*AbstractUsecase
	outputPortAnyLogic.AdminPermissionGroupLogic
}

func NewAdminPermissionGroupUsecase(oAbstractUsecase *AbstractUsecase, oAdminPermissionGroupLogic outputPortAnyLogic.AdminPermissionGroupLogic) usecasePortAnyLogic.AdminPermissionGroupUsecase {
	return &AdminPermissionGroupUsecase{
		AbstractUsecase:           oAbstractUsecase,
		AdminPermissionGroupLogic: oAdminPermissionGroupLogic,
	}
}

func (oSelf *AdminPermissionGroupUsecase) ShowTree() ([]*domain.AdminPermissionGroup, error) {

	aTree, oErr := oSelf.AdminPermissionGroupLogic.ShowTree()

	return aTree, oErr
}

func (oSelf *AdminPermissionGroupUsecase) ShowAdminPermissionGroupById(iId uint64) (*domain.AdminPermissionGroup, error) {

	oAdminPermissionGroup, oErr := oSelf.AdminPermissionGroupLogic.ShowAdminPermissionGroupById(iId)

	return oAdminPermissionGroup, oErr
}

func (oSelf *AdminPermissionGroupUsecase) ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, uint64, error) {

	aAdminPermissionGroups, iTotal, oErr := oSelf.AdminPermissionGroupLogic.ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aAdminPermissionGroups, uint64(iTotal), oErr
}
