package usecaseApplicationAnyAdminResource

import (
	bootstrap "example/bootstrap"
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type AdminUserUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	AdminUserModel outputPortAnyModel.AdminUserModel
	AdminUserLogic outputPortAnyLogic.AdminUserLogic
}

func NewAdminUserUsecase(oAdminUserModel outputPortAnyModel.AdminUserModel, oAdminUserLogic outputPortAnyLogic.AdminUserLogic, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminResource.AdminUserUsecase {
	return &AdminUserUsecase{
		AbstractUsecase: oAbstractUsecase,
		AdminUserModel:  oAdminUserModel,
		AdminUserLogic:  oAdminUserLogic,
	}
}

func (oSelf *AdminUserUsecase) AddOne(oAdminUser *domain.AdminUserValue) (bool, error) {

	oErr := oSelf.ValidatorHelper.Valiate(oAdminUser)

	if oErr != nil {
		return false, oErr
	}

	oErr = oSelf.AdminUserLogic.AddAminUser(oAdminUser)

	return true, oErr
}

func (oSelf *AdminUserUsecase) EditOne(oAdminUser *domain.AdminUserValue, iId uint) (bool, error) {

	// oAdminUser.  相當 *oAdminUser.
	// *oAdminUser.Name 相當 (*(oAdminUser).Name)

	if (*oAdminUser).Name != nil && *oAdminUser.Name != "" {
		return false, pkgUtility.NewDefaultError("name can not be modified", -1, 200)

	}

	oErr := oSelf.ValidatorHelper.Valiate(oAdminUser)

	if oErr != nil {
		return false, oErr
	}

	if (*oAdminUser).Password != nil && *(*oAdminUser).Password != "" {
		*oAdminUser.Password = pkgUtility.Md5(*oAdminUser.Password + bootstrap.CONFIG.TABLE.ADMIN_USER.PASSWORD)

	}

	oErr = oSelf.AdminUserLogic.EditAdminUserById(oAdminUser, iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AdminUserUsecase) RemoveOne(iId uint) (bool, error) {

	if iId == 1 {
		return false, pkgUtility.NewDefaultError("AdminUser.Id=1 cannot be deleted", -1, 200)

	}

	oErr := oSelf.AdminUserModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AdminUserUsecase) ShowOne(iId uint) (*domain.AdminUser, error) {

	oAdminUser, oErr := oSelf.AdminUserModel.ShowOneById(iId)

	return oAdminUser, oErr
}

func (oSelf *AdminUserUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint64, error) {

	aAdminUsers, iTotal, oErr := oSelf.AdminUserLogic.ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aAdminUsers, uint64(iTotal), oErr
}
