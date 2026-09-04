package usecaseApplicationAnyAdminResource

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type AdminUserUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	outputPortAnyModel.AdminUserModel
}

func NewAdminUserUsecase(oAdminUserModel outputPortAnyModel.AdminUserModel, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminResource.AdminUserUsecase {
	return &AdminUserUsecase{
		AbstractUsecase: oAbstractUsecase,
		AdminUserModel:  oAdminUserModel,
	}
}

func (oSelf *AdminUserUsecase) EditOne(oAdminUser *domain.AdminUserValue, iId uint) (bool, error) {

	_, oErr := oSelf.AdminUserModel.EditOneById(oAdminUser, iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AdminUserUsecase) RemoveOne(iId uint) (bool, error) {

	if iId == 1 {
		return false, pkgUtility.NewDefaultError("AdminUser.Id=1 無法刪除", -2, 200)

	}

	_, oErr := oSelf.AdminUserModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AdminUserUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint64, error) {

	aAdminUsers, oErr := oSelf.AdminUserModel.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, 0, oErr
	}

	iTotal, oErr := oSelf.AdminUserModel.TotalByFilters(aFilters)

	return aAdminUsers, iTotal, oErr
}
