package usecaseApplicationAnyModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.AdminPermissionGroupModel
}

func NewAdminPermissionGroupUsecase(oAdminPermissionGroupModel outputPortAnyModel.AdminPermissionGroupModel, oAbstractUsecase *AbstractUsecase) usecasePortAnyModel.AdminPermissionGroupUsecase {
	return &AdminPermissionGroupUsecase{
		AbstractUsecase:           oAbstractUsecase,
		AdminPermissionGroupModel: oAdminPermissionGroupModel,
	}
}

func (oSelf *AdminPermissionGroupUsecase) AddOne(oValue *domain.AdminPermissionGroupValue) error {

	oErr := oSelf.AdminPermissionGroupModel.AddOne(oValue)

	return oErr
}

func (oSelf *AdminPermissionGroupUsecase) ShowOneById(iId uint64) (*domain.AdminPermissionGroup, error) {

	oAdminPermissionGroup, oErr := oSelf.AdminPermissionGroupModel.ShowOneById(iId)

	return oAdminPermissionGroup, oErr
}

func (oSelf *AdminPermissionGroupUsecase) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, error) {

	aAdminPermissionGroups, oErr := oSelf.AdminPermissionGroupModel.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aAdminPermissionGroups, oErr
}

func (oSelf *AdminPermissionGroupUsecase) EditOneById(oValue *domain.AdminPermissionGroupValue, iId uint64) error {

	oErr := oSelf.AdminPermissionGroupModel.EditOneById(oValue, iId)

	return oErr
}

func (oSelf *AdminPermissionGroupUsecase) RemoveOneById(iId uint64) error {

	oErr := oSelf.AdminPermissionGroupModel.RemoveOneById(iId)

	return oErr
}

func (oSelf *AdminPermissionGroupUsecase) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	iTotal, oErr := oSelf.AdminPermissionGroupModel.TotalByFilters(aFilters)

	return uint64(iTotal), oErr
}
