package usecaseApplicationAnyAdminResource

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkgInput "example/pkg/input"
)

type AdminRoleUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	outputPortAnyModel.AdminRoleModel
}

func NewAdminRoleUsecase(oAdminRoleModel outputPortAnyModel.AdminRoleModel, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminResource.AdminRoleUsecase {
	return &AdminRoleUsecase{
		AbstractUsecase: oAbstractUsecase,
		AdminRoleModel:  oAdminRoleModel,
	}
}

func (oSelf *AdminRoleUsecase) AddOne(oAdminRole *domain.AdminRoleValue) error {

	return oSelf.AdminRoleModel.AddOne(oAdminRole)
}

func (oSelf *AdminRoleUsecase) ShowOne(iId uint) (*domain.AdminRole, error) {

	oAdminRole, oErr := oSelf.AdminRoleModel.ShowOneById(iId)

	return oAdminRole, oErr
}

func (oSelf *AdminRoleUsecase) EditOne(oAdminRole *domain.AdminRoleValue, iId uint) error {

	return oSelf.AdminRoleModel.EditOneById(oAdminRole, iId)
}

func (oSelf *AdminRoleUsecase) RemoveOne(iId uint) error {

	return oSelf.AdminRoleModel.RemoveOneById(iId)
}

func (oSelf *AdminRoleUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminRole, uint, error) {

	aAdminRoles, oErr := oSelf.AdminRoleModel.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, 0, oErr
	}

	iTotal, oErr := oSelf.AdminRoleModel.TotalByFilters(aFilters)

	return aAdminRoles, iTotal, oErr
}
