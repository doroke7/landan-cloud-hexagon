package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
	pkgInput "example/pkg/input"
)

type AdminUserLogic struct {
	*resourceBase.AbstractResource
}

func NewAdminUserLogic(oAbstractLogic *resourceBase.AbstractResource) outputPortAnyLogic.AdminUserLogic {
	return &AdminUserLogic{
		AbstractResource: oAbstractLogic,
	}
}

func (oSelf *AdminUserLogic) ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint64, error) {

	oRequest := &pbResourceLogic.AdminUserShowAdminUsersTotalByFiltersWithSortersPaginationInput{
		Filters: oSelf.ToFilters(aFilters),
		Sorters: oSelf.ToSorters(aSorters),
	}

	if oPagination != nil {
		oRequest.Pagination = oSelf.ToPagination(oPagination)
	}

	oResponse, oErr := oSelf.ResourceLogicClient.AdminUser.ShowAdminUsersTotalByFiltersWithSortersPagination(oSelf.Context, oRequest)

	aAdminUsers := make([]*domain.AdminUser, 0, len(oResponse.GetAdminUsers()))
	for _, oOne := range oResponse.GetAdminUsers() {
		aAdminUsers = append(aAdminUsers, &domain.AdminUser{
			Id:        uint(oOne.GetId()),
			Name:      oOne.GetName(),
			Password:  oOne.GetPassword(),
			CreatedAt: oOne.GetCreatedAt().AsTime(),
			UpdatedAt: oOne.GetUpdatedAt().AsTime(),
			DeletedAt: oOne.GetDeletedAt().AsTime(),
		})
	}

	iTotal := oResponse.GetTotal()
	return aAdminUsers, iTotal, oErr
}

func (oSelf *AdminUserLogic) AddAminUser(oValue *domain.AdminUserValue) error {

	aAdminRoleIds := make([]uint32, 0, len(oValue.AdminRoleIds))
	for _, iAdminRoleId := range oValue.AdminRoleIds {
		aAdminRoleIds = append(aAdminRoleIds, uint32(iAdminRoleId))
	}

	oRequest := &pbResourceLogic.AdminUserAddAminUserInput{
		Variable: &pbResourceLogic.AdminUserVariable{
			Name:         oValue.Name,
			Password:     oValue.Password,
			AdminRoleIds: aAdminRoleIds,
		},
	}

	_, oErr := oSelf.ResourceLogicClient.AdminUser.AddAminUser(oSelf.Context, oRequest)

	return oErr
}
