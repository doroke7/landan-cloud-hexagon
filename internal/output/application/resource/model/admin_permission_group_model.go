package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
	pkgProtoToDomain "example/pkg/proto_to_domain"
)

type AdminPermissionGroupModel struct {
	*AbstractModel
}

func NewAdminPermissionGroupModel(oAbstractModel *AbstractModel) outputPortAnyModel.AdminPermissionGroupModel {
	oModel := &AdminPermissionGroupModel{
		AbstractModel: oAbstractModel,
	}

	return oModel
}

func (oSelf *AdminPermissionGroupModel) AddOne(oValue *domain.AdminPermissionGroupVariable) error {

	oRequest := &pbResourceModel.AdminPermissionGroupAddOneInput{
		Variable: &pbResource.AdminPermissionGroupVariable{
			Key:  oValue.Key,
			Name: oValue.Name,
		},
	}

	_, oErr := oSelf.ResourceModelClient.AdminPermissionGroup.AddOne(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminPermissionGroupModel) ShowOnesByParentId(iParentId uint64) ([]*domain.AdminPermissionGroup, error) {

	oRequest := &pbResourceModel.AdminPermissionGroupShowOnesByParentIdInput{ParentId: iParentId}
	oResponse, oErr := oSelf.ResourceModelClient.AdminPermissionGroup.ShowOnesByParentId(oSelf.Context, oRequest)

	if oErr != nil {
		return nil, oErr
	}

	aAdminPermissionGroups := make([]*domain.AdminPermissionGroup, 0, len(oResponse.GetAdminPermissionGroups()))
	for _, oProto := range oResponse.GetAdminPermissionGroups() {
		oAdminPermissionGroup := pkgProtoToDomain.AdminPermissionGroup(oProto)
		aAdminPermissionGroups = append(aAdminPermissionGroups, &oAdminPermissionGroup)
	}

	return aAdminPermissionGroups, nil
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
		oAdminPermissionGroup := pkgProtoToDomain.AdminPermissionGroup(oProto)
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

func (oSelf *AdminPermissionGroupModel) EditOneById(oValue *domain.AdminPermissionGroupVariable, iId uint64) error {

	oRequest := &pbResourceModel.AdminPermissionGroupEditOneByIdInput{
		Id: iId,
		Variable: &pbResource.AdminPermissionGroupVariable{
			Key:  oValue.Key,
			Name: oValue.Name,
		},
	}

	_, oErr := oSelf.ResourceModelClient.AdminPermissionGroup.EditOneById(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminPermissionGroupModel) RemoveOneById(iId uint64) error {

	oRequest := &pbResourceModel.AdminPermissionGroupRemoveOneByIdInput{Id: iId}
	_, oErr := oSelf.ResourceModelClient.AdminPermissionGroup.RemoveOneById(oSelf.Context, oRequest)

	return oErr
}
