package usecaseApplicationAnyLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pkgInput "example/pkg/input"
)

type AdminUserUsecase struct {
	*AbstractUsecase
	outputPortAnyLogic.AdminUserLogic
}

func NewAdminUserUsecase(oAbstractUsecase *AbstractUsecase, oAdminUserLogic outputPortAnyLogic.AdminUserLogic) usecasePortAnyLogic.AdminUserUsecase {
	return &AdminUserUsecase{
		AbstractUsecase: oAbstractUsecase,
		AdminUserLogic:  oAdminUserLogic,
	}
}

func (oSelf *AdminUserUsecase) ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint64, error) {

	aAdminUsers, iTotal, oErr := oSelf.AdminUserLogic.ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aAdminUsers, uint64(iTotal), oErr
}

func (oSelf *AdminUserUsecase) ShowAdminUserById(iId uint64) (*domain.AdminUser, error) {

	oAdminUser, oErr := oSelf.AdminUserLogic.ShowAdminUserById(iId)

	return oAdminUser, oErr
}

func (oSelf *AdminUserUsecase) AddAminUser(oValue *domain.AdminUserVariable) error {

	oErr := oSelf.AdminUserLogic.AddAdminUser(oValue)

	return oErr
}

func (oSelf *AdminUserUsecase) EditAdminUserById(oValue *domain.AdminUserVariable, iId uint64) error {

	oErr := oSelf.AdminUserLogic.EditAdminUserById(oValue, iId)

	return oErr
}

func (oSelf *AdminUserUsecase) RemoveAdminUserById(iId uint64) error {

	oErr := oSelf.AdminUserLogic.RemoveAdminUserById(iId)

	return oErr
}
