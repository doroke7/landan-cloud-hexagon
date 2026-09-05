package usecaseApplicationAnyModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
)

type AdminRoleUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.AdminRoleModel
}

func NewAdminRoleUsecase(oAdminRoleModel outputPortAnyModel.AdminRoleModel, oAbstractUsecase *AbstractUsecase) usecasePortAnyModel.AdminRoleUsecase {
	return &AdminRoleUsecase{
		AbstractUsecase: oAbstractUsecase,
		AdminRoleModel:  oAdminRoleModel,
	}
}

func (oSelf *AdminRoleUsecase) AddOne(oAdminRole *domain.AdminRoleValue) error {

	return oSelf.AdminRoleModel.AddOne(oAdminRole)
}

func (oSelf *AdminRoleUsecase) EditOneById(oAdminRole *domain.AdminRoleValue, iId uint64) error {

	return oSelf.AdminRoleModel.EditOneById(oAdminRole, uint(iId))
}

func (oSelf *AdminRoleUsecase) RemoveOneById(iId uint64) error {

	return oSelf.AdminRoleModel.RemoveOneById(uint(iId))
}

func (oSelf *AdminRoleUsecase) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	iTotal, oErr := oSelf.AdminRoleModel.TotalByFilters(aFilters)

	return uint64(iTotal), oErr
}
