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

func (oSelf *AdminUserUsecase) AddAminUser(oValue *domain.AdminUserValue) error {

	oErr := oSelf.AdminUserLogic.AddAminUser(oValue)

	return oErr
}

func (oSelf *AdminUserUsecase) EditAdminUserById(oValue *domain.AdminUserValue, iId uint64) error {

	oErr := oSelf.AdminUserLogic.EditAdminUserById(oValue, uint(iId))

	return oErr
}
