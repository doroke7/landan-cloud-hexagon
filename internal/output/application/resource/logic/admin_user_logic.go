package outputApplicationResourceLogic

import (
	"errors"

	domain "example/internal/domain"
	outputApplicationResource "example/internal/output/application/resource"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResource "example/pb/resource"
	pbResourceLogic "example/pb/resource/logic"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
)

func protoAdminRoleToDomainAdminRole(oProto *pbResource.AdminRole) domain.AdminRole {
	if oProto == nil {
		return domain.AdminRole{}
	}

	return domain.AdminRole{
		Id:        uint(oProto.GetId()),
		Key:       oProto.GetKey(),
		Name:      oProto.GetName(),
		CreatedAt: oProto.GetCreatedAt().AsTime(),
		UpdatedAt: oProto.GetUpdatedAt().AsTime(),
		DeletedAt: oProto.GetDeletedAt().AsTime(),
	}
}

type AdminUserLogic struct {
	*outputApplicationResource.AbstractResource
}

func NewAdminUserLogic(oAbstractLogic *outputApplicationResource.AbstractResource) outputPortAnyLogic.AdminUserLogic {
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
		aAdminRoles := make([]domain.AdminRole, 0, len(oOne.GetAdminRoles()))
		for _, oProtoAdminRole := range oOne.GetAdminRoles() {
			aAdminRoles = append(aAdminRoles, protoAdminRoleToDomainAdminRole(oProtoAdminRole))
		}

		aAdminUsers = append(aAdminUsers, &domain.AdminUser{
			Id:         uint(oOne.GetId()),
			Name:       oOne.GetName(),
			Password:   oOne.GetPassword(),
			CreatedAt:  oOne.GetCreatedAt().AsTime(),
			UpdatedAt:  oOne.GetUpdatedAt().AsTime(),
			DeletedAt:  oOne.GetDeletedAt().AsTime(),
			AdminRoles: aAdminRoles,
		})
	}

	iTotal := uint64(oResponse.GetTotal())
	return aAdminUsers, iTotal, oErr
}

func (oSelf *AdminUserLogic) AddAminUser(oValue *domain.AdminUserValue) error {

	aAdminRoleIds := make([]uint64, 0, len(oValue.AdminRoleIds))
	for _, iAdminRoleId := range oValue.AdminRoleIds {
		aAdminRoleIds = append(aAdminRoleIds, uint64(iAdminRoleId))
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

func (oSelf *AdminUserLogic) EditAdminUserById(oValue *domain.AdminUserValue, iId uint64) error {

	aAdminRoleIds := make([]uint64, 0, len(oValue.AdminRoleIds))
	for _, iAdminRoleId := range oValue.AdminRoleIds {
		aAdminRoleIds = append(aAdminRoleIds, uint64(iAdminRoleId))
	}

	oRequest := &pbResourceLogic.AdminUserEditAdminUserByIdInput{
		Variable: &pbResourceLogic.AdminUserVariable{
			Name:         oValue.Name,
			Password:     oValue.Password,
			AdminRoleIds: aAdminRoleIds,
		},
		Id: uint64(iId),
	}

	_, oErr := oSelf.ResourceLogicClient.AdminUser.EditAdminUserById(oSelf.Context, oRequest)

	return oErr
}

// ShowAdminUserById proto 沒有對應 rpc，改用既有的 list rpc 帶 id 過濾、size=1 取第一筆。
func (oSelf *AdminUserLogic) ShowAdminUserById(iId uint64) (*domain.AdminUser, error) {
	sField := "id"
	sOperator := "eq"
	iSize := uint(1)
	iPage := uint(1)

	aFilters := []*pkgInput.Filter{
		{Field: &sField, Operator: &sOperator, Value: float64(iId)},
	}
	oPagination := &pkgInput.Pagination{Size: &iSize, Page: &iPage}

	aAdminUsers, _, oErr := oSelf.ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters, nil, oPagination)
	if oErr != nil {
		return nil, oErr
	}

	if len(aAdminUsers) == 0 {
		return nil, errors.New("record not found")
	}

	oAdminUser := aAdminUsers[0]

	return oAdminUser, nil
}

// RemoveAdminUserById logic proto 沒有 remove rpc，改打 model client 的 AdminUser.RemoveOneById。
func (oSelf *AdminUserLogic) RemoveAdminUserById(iId uint64) error {
	oRequest := &pbResourceModel.AdminUserRemoveOneByIdInput{Id: uint32(iId)}

	_, oErr := oSelf.ResourceModelClient.AdminUser.RemoveOneById(oSelf.Context, oRequest)

	return oErr
}
