package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
)

func protoAdminPermissionGroupToDomainAdminPermissionGroup(oProto *pbResource.AdminPermissionGroup) domain.AdminPermissionGroup {
	if oProto == nil {
		return domain.AdminPermissionGroup{}
	}

	oAdminPermissionGroup := domain.AdminPermissionGroup{
		Id:        uint64(oProto.GetId()),
		ParentId:  uint64(oProto.GetParentId()),
		Key:       oProto.GetKey(),
		Name:      oProto.GetName(),
		CreatedAt: oProto.GetCreatedAt().AsTime(),
		UpdatedAt: oProto.GetUpdatedAt().AsTime(),
		DeletedAt: oProto.GetDeletedAt().AsTime(),
	}

	if oParent := oProto.GetParent(); oParent != nil {
		oParentDomain := protoAdminPermissionGroupToDomainAdminPermissionGroup(oParent)
		oAdminPermissionGroup.Parent = &oParentDomain
	}

	for _, oChild := range oProto.GetChildren() {
		oAdminPermissionGroup.Children = append(oAdminPermissionGroup.Children, protoAdminPermissionGroupToDomainAdminPermissionGroup(oChild))
	}

	return oAdminPermissionGroup
}

type AdminPermissionGroupModel struct {
	*AbstractModel
}

func NewAdminPermissionGroupModel(oAbstractModel *AbstractModel) outputPortAnyModel.AdminPermissionGroupModel {
	return &AdminPermissionGroupModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *AdminPermissionGroupModel) AddOne(oValue *domain.AdminPermissionGroupValue) error {

	oRequest := &pbResourceModel.AdminPermissionGroupAddOneInput{
		Variable: &pbResourceModel.AdminPermissionGroupVariable{
			Key:  oValue.Key,
			Name: oValue.Name,
		},
	}

	_, oErr := oSelf.ResourceModelClient.AdminPermissionGroup.AddOne(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminPermissionGroupModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, error) {

	oRequest := &pbResourceModel.AdminPermissionGroupShowOnesByFiltersWithSortersPaginationInput{
		Filters:    oSelf.ToFilters(aFilters),
		Sorters:    oSelf.ToSorters(aSorters),
		Pagination: oSelf.ToPagination(oPagination),
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminPermissionGroup.ShowOnesByFiltersWithSortersPagination(oSelf.Context, oRequest)

	if oErr != nil {
		return nil, oErr
	}

	aAdminPermissionGroups := make([]*domain.AdminPermissionGroup, 0, len(oResponse.GetAdminPermissionGroups()))
	for _, oProto := range oResponse.GetAdminPermissionGroups() {
		oAdminPermissionGroup := protoAdminPermissionGroupToDomainAdminPermissionGroup(oProto)
		aAdminPermissionGroups = append(aAdminPermissionGroups, &oAdminPermissionGroup)
	}

	return aAdminPermissionGroups, nil
}

func (oSelf *AdminPermissionGroupModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.AdminPermissionGroupTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminPermissionGroup.TotalByFilters(oSelf.Context, oRequest)

	iTotal := uint64(oResponse.GetTotal())

	return iTotal, oErr
}

func (oSelf *AdminPermissionGroupModel) EditOneById(oValue *domain.AdminPermissionGroupValue, iId uint64) error {

	oRequest := &pbResourceModel.AdminPermissionGroupEditOneByIdInput{
		Id: iId,
		Variable: &pbResourceModel.AdminPermissionGroupVariable{
			Key:  oValue.Key,
			Name: oValue.Name,
		},
	}

	_, oErr := oSelf.ResourceModelClient.AdminPermissionGroup.EditOneById(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminPermissionGroupModel) RemoveOneById(iId uint64) error {

	_, oErr := oSelf.ResourceModelClient.AdminPermissionGroup.RemoveOneById(
		oSelf.Context,
		&pbResourceModel.AdminPermissionGroupRemoveOneByIdInput{Id: iId},
	)

	return oErr
}
