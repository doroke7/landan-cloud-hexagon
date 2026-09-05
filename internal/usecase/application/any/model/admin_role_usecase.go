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

func (oSelf *AdminRoleUsecase) AddOne(oAdminRole *domain.AdminRoleValue) (bool, error) {

	oErr := oSelf.AdminRoleModel.AddOne(oAdminRole)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AdminRoleUsecase) EditOneById(oAdminRole *domain.AdminRoleValue, iId uint) (bool, error) {

	oErr := oSelf.AdminRoleModel.EditOneById(oAdminRole, iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AdminRoleUsecase) RemoveOneById(iId uint) (bool, error) {

	oErr := oSelf.AdminRoleModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AdminRoleUsecase) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	iTotal, oErr := oSelf.AdminRoleModel.TotalByFilters(aFilters)

	return uint64(iTotal), oErr
}
