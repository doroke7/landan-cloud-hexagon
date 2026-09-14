package outputApplicationResourceLogic

import (
	"errors"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResource "example/pb/resource"
	pbResourceLogic "example/pb/resource/logic"
	pkgInput "example/pkg/input"
)

func protoAdminRoleToDomainAdminRole(oProto *pbResource.AdminRole) domain.AdminRole {
	if oProto == nil {
		return domain.AdminRole{}
	}

	oAdminRole := domain.AdminRole{
		Id:        uint64(oProto.GetId()),
		Key:       oProto.GetKey(),
		Name:      oProto.GetName(),
		CreatedAt: oProto.GetCreatedAt().AsTime(),
		UpdatedAt: oProto.GetUpdatedAt().AsTime(),
		DeletedAt: oProto.GetDeletedAt().AsTime(),
	}

	return oAdminRole
}

type AdminUserLogic struct {
	*AbstractLogic
}

func NewAdminUserLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminUserLogic {
	oLogic := &AdminUserLogic{
		AbstractLogic: oAbstractLogic,
	}

	return oLogic
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

		oAdminUser := &domain.AdminUser{
			Id:         uint64(oOne.GetId()),
			Name:       oOne.GetName(),
			Password:   oOne.GetPassword(),
			CreatedAt:  oOne.GetCreatedAt().AsTime(),
			UpdatedAt:  oOne.GetUpdatedAt().AsTime(),
			DeletedAt:  oOne.GetDeletedAt().AsTime(),
			AdminRoles: aAdminRoles,
		}
		aAdminUsers = append(aAdminUsers, oAdminUser)
	}

	iTotal := uint64(oResponse.GetTotal())
	return aAdminUsers, iTotal, oErr
}

func (oSelf *AdminUserLogic) AddAminUser(oValue *domain.AdminUserVariable) error {

	aAdminRoleIds := make([]uint64, 0, len(oValue.AdminRoleIds))
	for _, iAdminRoleId := range oValue.AdminRoleIds {
		aAdminRoleIds = append(aAdminRoleIds, uint64(iAdminRoleId))
	}

	oRequest := &pbResourceLogic.AdminUserAddAminUserInput{
		Value: &pbResource.AdminUserValue{
			Name:         oValue.Name,
			Password:     oValue.Password,
			AdminRoleIds: aAdminRoleIds,
		},
	}

	_, oErr := oSelf.ResourceLogicClient.AdminUser.AddAminUser(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminUserLogic) EditAdminUserById(oValue *domain.AdminUserVariable, iId uint64) error {

	aAdminRoleIds := make([]uint64, 0, len(oValue.AdminRoleIds))
	for _, iAdminRoleId := range oValue.AdminRoleIds {
		aAdminRoleIds = append(aAdminRoleIds, uint64(iAdminRoleId))
	}

	oRequest := &pbResourceLogic.AdminUserEditAdminUserByIdInput{
		Value: &pbResource.AdminUserValue{
			Name:         oValue.Name,
			Password:     oValue.Password,
			AdminRoleIds: aAdminRoleIds,
		},
		Id: uint64(iId),
	}

	_, oErr := oSelf.ResourceLogicClient.AdminUser.EditAdminUserById(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminUserLogic) ShowAdminUserById(iId uint64) (*domain.AdminUser, error) {
	oRequest := &pbResourceLogic.AdminUserShowAdminUserByIdInput{Id: iId}

	oResponse, oErr := oSelf.ResourceLogicClient.AdminUser.ShowAdminUserById(oSelf.Context, oRequest)
	if oErr != nil {
		return nil, oErr
	}

	if len(oResponse.GetAdminUsers()) == 0 {
		return nil, errors.New("record not found")
	}

	oOne := oResponse.GetAdminUsers()[0]

	aAdminRoles := make([]domain.AdminRole, 0, len(oOne.GetAdminRoles()))
	for _, oProtoAdminRole := range oOne.GetAdminRoles() {
		aAdminRoles = append(aAdminRoles, protoAdminRoleToDomainAdminRole(oProtoAdminRole))
	}

	oAdminUser := &domain.AdminUser{
		Id:         uint64(oOne.GetId()),
		Name:       oOne.GetName(),
		Password:   oOne.GetPassword(),
		CreatedAt:  oOne.GetCreatedAt().AsTime(),
		UpdatedAt:  oOne.GetUpdatedAt().AsTime(),
		DeletedAt:  oOne.GetDeletedAt().AsTime(),
		AdminRoles: aAdminRoles,
	}

	return oAdminUser, nil
}

func (oSelf *AdminUserLogic) RemoveAdminUserById(iId uint64) error {
	oRequest := &pbResourceLogic.AdminUserRemoveAdminUserByIdInput{Id: iId}

	_, oErr := oSelf.ResourceLogicClient.AdminUser.RemoveAdminUserById(oSelf.Context, oRequest)

	return oErr
}
