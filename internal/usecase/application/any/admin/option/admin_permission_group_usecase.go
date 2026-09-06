package usecaseApplicationAnyAdminOption

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminOption "example/internal/usecase/port/any/admin/option"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	AdminPermissionGroupLogic outputPortAnyLogic.AdminPermissionGroupLogic
}

func NewAdminPermissionGroupUsecase(oAdminPermissionGroupLogic outputPortAnyLogic.AdminPermissionGroupLogic, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminOption.AdminPermissionGroupUsecase {
	return &AdminPermissionGroupUsecase{
		AbstractUsecase:           oAbstractUsecase,
		AdminPermissionGroupLogic: oAdminPermissionGroupLogic,
	}
}

func (oSelf *AdminPermissionGroupUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, uint64, error) {

	aAdminPermissionGroups, oErr := oSelf.AdminPermissionGroupLogic.ShowAdminPermissionGroups()
	if oErr != nil {
		return nil, 0, oErr
	}

	iTotal := len(aAdminPermissionGroups)

	return aAdminPermissionGroups, uint64(iTotal), oErr
}

func (oSelf *AdminPermissionGroupUsecase) ShowTree() ([]*domain.AdminPermissionGroup, error) {

	aTree, oErr := oSelf.AdminPermissionGroupLogic.ShowTree()

	return aTree, oErr
}
