package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
)

type AdminRoleModel struct {
	*resourceBase.AbstractResource
}

func NewAdminRoleModel(oAbstractModel *resourceBase.AbstractResource) outputPortAnyModel.AdminRoleModel {
	return &AdminRoleModel{
		AbstractResource: oAbstractModel,
	}
}

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

func domainAdminRoleValueToProtoAdminRoleVariable(oValue *domain.AdminRoleValue) *pbResourceModel.AdminRoleVariable {
	return &pbResourceModel.AdminRoleVariable{
		Key:  oValue.Key,
		Name: oValue.Name,
	}
}

func (oSelf *AdminRoleModel) AddOne(oAdminRole *domain.AdminRoleValue) error {

	oRequest := &pbResourceModel.AdminRoleAddOneInput{
		Variable: domainAdminRoleValueToProtoAdminRoleVariable(oAdminRole),
	}

	_, oErr := oSelf.ResourceModelClient.AdminRole.AddOne(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminRoleModel) ShowOneById(iId uint64) (*domain.AdminRole, error) {

	oResponse, oErr := oSelf.ResourceModelClient.AdminRole.ShowOneById(
		oSelf.Context,
		&pbResourceModel.AdminRoleShowOneByIdInput{Id: uint64(iId)},
	)

	if oErr != nil {
		return nil, oErr
	}

	oProtoAdminRole := oResponse.GetAdminRole()
	if oProtoAdminRole.GetId() == 0 {
		return nil, nil
	}

	oAdminRole := protoAdminRoleToDomainAdminRole(oProtoAdminRole)

	return &oAdminRole, nil
}

func (oSelf *AdminRoleModel) EditOneById(oAdminRole *domain.AdminRoleValue, iId uint64) error {

	oRequest := &pbResourceModel.AdminRoleEditOneByIdInput{
		Id:       uint64(iId),
		Variable: domainAdminRoleValueToProtoAdminRoleVariable(oAdminRole),
	}

	_, oErr := oSelf.ResourceModelClient.AdminRole.EditOneById(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminRoleModel) RemoveOneById(iId uint64) error {

	_, oErr := oSelf.ResourceModelClient.AdminRole.RemoveOneById(
		oSelf.Context,
		&pbResourceModel.AdminRoleRemoveOneByIdInput{Id: uint64(iId)},
	)

	return oErr
}

func (oSelf *AdminRoleModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminRole, error) {

	oRequest := &pbResourceModel.AdminRoleShowOnesByFiltersWithSortersPaginationInput{
		Filters:    oSelf.ToFilters(aFilters),
		Sorters:    oSelf.ToSorters(aSorters),
		Pagination: oSelf.ToPagination(oPagination),
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminRole.ShowOnesByFiltersWithSortersPagination(oSelf.Context, oRequest)

	if oErr != nil {
		return nil, oErr
	}

	aAdminRoles := make([]*domain.AdminRole, 0, len(oResponse.GetAdminRoles()))
	for _, oOne := range oResponse.GetAdminRoles() {
		oAdminRole := protoAdminRoleToDomainAdminRole(oOne)
		aAdminRoles = append(aAdminRoles, &oAdminRole)
	}

	return aAdminRoles, nil
}

func (oSelf *AdminRoleModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.AdminRoleTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminRole.TotalByFilters(oSelf.Context, oRequest)

	iTotal := uint64(oResponse.GetTotal())

	return iTotal, oErr
}
