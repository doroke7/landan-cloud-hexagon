package usecaseApplicationAnyAdminResource

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkgInput "example/pkg/input"
)

type AdminRoleUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	AdminRoleModel outputPortAnyModel.AdminRoleModel
	AdminRoleLogic outputPortAnyLogic.AdminRoleLogic
}

func NewAdminRoleUsecase(oAdminRoleModel outputPortAnyModel.AdminRoleModel, oAdminRoleLogic outputPortAnyLogic.AdminRoleLogic, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminResource.AdminRoleUsecase {
	return &AdminRoleUsecase{
		AbstractUsecase: oAbstractUsecase,
		AdminRoleModel:  oAdminRoleModel,
		AdminRoleLogic:  oAdminRoleLogic,
	}
}

func (oSelf *AdminRoleUsecase) AddOne(oAdminRole *domain.AdminRoleVariable) error {

	oErr := oSelf.AdminRoleLogic.AddAdminRole(oAdminRole)

	return oErr
}

func (oSelf *AdminRoleUsecase) ShowOne(iId uint64) (*domain.AdminRole, error) {

	oAdminRole, oErr := oSelf.AdminRoleModel.ShowOneById(iId)

	return oAdminRole, oErr
}

func (oSelf *AdminRoleUsecase) EditOne(oAdminRole *domain.AdminRoleVariable, iId uint64) error {

	oErr := oSelf.AdminRoleModel.EditOneById(oAdminRole, iId)

	return oErr
}

func (oSelf *AdminRoleUsecase) RemoveOne(iId uint64) error {

	oErr := oSelf.AdminRoleModel.RemoveOneById(iId)

	return oErr
}

func (oSelf *AdminRoleUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminRole, uint64, error) {

	aAdminRoles, oErr := oSelf.AdminRoleModel.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, 0, oErr
	}

	iTotal, oErr := oSelf.AdminRoleModel.TotalByFilters(aFilters)

	return aAdminRoles, uint64(iTotal), oErr
}
