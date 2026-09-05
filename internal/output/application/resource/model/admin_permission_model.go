package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
)

type AdminPermissionModel struct {
	*resourceBase.AbstractResource
}

func NewAdminPermissionModel(oAbstractModel *resourceBase.AbstractResource) outputPortAnyModel.AdminPermissionModel {
	return &AdminPermissionModel{
		AbstractResource: oAbstractModel,
	}
}

func protoAdminPermissionToDomainAdminPermission(oProto *pbResource.AdminPermission) domain.AdminPermission {
	if oProto == nil {
		return domain.AdminPermission{}
	}

	return domain.AdminPermission{
		Id:        uint(oProto.GetId()),
		Type:      uint8(oProto.GetType()),
		Key:       oProto.GetKey(),
		Name:      oProto.GetName(),
		CreatedAt: oProto.GetCreatedAt().AsTime(),
		UpdatedAt: oProto.GetUpdatedAt().AsTime(),
		DeletedAt: oProto.GetDeletedAt().AsTime(),
	}
}

func domainAdminPermissionValueToProtoAdminPermissionVariable(oValue *domain.AdminPermissionValue) *pbResourceModel.AdminPermissionVariable {
	oVariable := &pbResourceModel.AdminPermissionVariable{
		Key:  oValue.Key,
		Name: oValue.Name,
	}

	if oValue.Type != nil {
		iType := uint32(*oValue.Type)
		oVariable.Type = &iType
	}

	return oVariable
}

func (oSelf *AdminPermissionModel) AddOne(oAdminPermission *domain.AdminPermissionValue) error {

	oRequest := &pbResourceModel.AdminPermissionAddOneInput{
		Variable: domainAdminPermissionValueToProtoAdminPermissionVariable(oAdminPermission),
	}

	_, oErr := oSelf.ResourceModelClient.AdminPermission.AddOne(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminPermissionModel) ShowOneById(iId uint) (*domain.AdminPermission, error) {

	oResponse, oErr := oSelf.ResourceModelClient.AdminPermission.ShowOneById(
		oSelf.Context,
		&pbResourceModel.AdminPermissionShowOneByIdInput{Id: uint32(iId)},
	)

	if oErr != nil {
		return nil, oErr
	}

	oProtoAdminPermission := oResponse.GetAdminPermission()
	if oProtoAdminPermission.GetId() == 0 {
		return nil, nil
	}

	oAdminPermission := protoAdminPermissionToDomainAdminPermission(oProtoAdminPermission)

	return &oAdminPermission, nil
}

func (oSelf *AdminPermissionModel) EditOneById(oAdminPermission *domain.AdminPermissionValue, iId uint) error {

	oRequest := &pbResourceModel.AdminPermissionEditOneByIdInput{
		Id:       uint32(iId),
		Variable: domainAdminPermissionValueToProtoAdminPermissionVariable(oAdminPermission),
	}

	_, oErr := oSelf.ResourceModelClient.AdminPermission.EditOneById(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminPermissionModel) RemoveOneById(iId uint) error {

	_, oErr := oSelf.ResourceModelClient.AdminPermission.RemoveOneById(
		oSelf.Context,
		&pbResourceModel.AdminPermissionRemoveOneByIdInput{Id: uint32(iId)},
	)

	return oErr
}

func (oSelf *AdminPermissionModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermission, error) {

	oRequest := &pbResourceModel.AdminPermissionShowOnesByFiltersWithSortersPaginationInput{
		Filters:    oSelf.ToFilters(aFilters),
		Sorters:    oSelf.ToSorters(aSorters),
		Pagination: oSelf.ToPagination(oPagination),
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminPermission.ShowOnesByFiltersWithSortersPagination(oSelf.Context, oRequest)

	if oErr != nil {
		return nil, oErr
	}

	aAdminPermissions := make([]*domain.AdminPermission, 0, len(oResponse.GetAdminPermissions()))
	for _, oOne := range oResponse.GetAdminPermissions() {
		oAdminPermission := protoAdminPermissionToDomainAdminPermission(oOne)
		aAdminPermissions = append(aAdminPermissions, &oAdminPermission)
	}

	return aAdminPermissions, nil
}

func (oSelf *AdminPermissionModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint, error) {

	oRequest := &pbResourceModel.AdminPermissionTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminPermission.TotalByFilters(oSelf.Context, oRequest)

	iTotal := uint(oResponse.GetTotal())

	return iTotal, oErr
}
