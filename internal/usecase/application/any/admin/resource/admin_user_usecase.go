package resource

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkg "example/pkg"
)

type AdminUserUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.AdminUserModel
}

func NewAdminUserUsecase(oAdminUserModel outputPortAnyModel.AdminUserModel, oAbstractUsecase *AbstractUsecase) usecasePortAnyAdminResource.AdminUserUsecase {
	return &AdminUserUsecase{
		AbstractUsecase: oAbstractUsecase,
		AdminUserModel:  oAdminUserModel,
	}
}

func (oSelf *AdminUserUsecase) ShowOnes(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.AdminUser, uint64, error) {

	aWheres := pkg.FiltersToMysqlWheres([]string{}, aFilters)
	aOrders := pkg.SortersToMysqlOrders([]string{}, aSorters)

	oLimit := pkg.PaginationToMysqlLimit(oPagination)

	aAdminUsers, oErr := oSelf.AdminUserModel.ShowOnesByWheresWithOrdersLimit(aWheres, aOrders, oLimit)
	if oErr != nil {
		return nil, 0, oErr
	}

	iTotal, oErr := oSelf.AdminUserModel.TotalByWheres(aWheres)

	return aAdminUsers, iTotal, oErr
}
