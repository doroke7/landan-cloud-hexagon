package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResourceModel "example/pb/resource/model"
	pkgDomainToProto "example/pkg/domain_to_proto"
	pkgInput "example/pkg/input"
	pkgProtoToDomain "example/pkg/proto_to_domain"
)

type AdminPermissionModel struct {
	*AbstractModel
}

func NewAdminPermissionModel(oAbstractModel *AbstractModel) outputPortAnyModel.AdminPermissionModel {
	oModel := &AdminPermissionModel{
		AbstractModel: oAbstractModel,
	}

	return oModel
}

func domainAdminPermissionValueToProtoAdminPermissionValue(oValue *domain.AdminPermissionValue) *pbResourceModel.AdminPermissionValue {
	oVariable := &pbResourceModel.AdminPermissionValue{
		Key:  oValue.Key,
		Name: oValue.Name,
	}

	if oValue.Type != nil {
		iType := uint64(*oValue.Type)
		oVariable.Type = &iType
	}

	return oVariable
}

func (oSelf *AdminPermissionModel) AddOne(oAdminPermission *domain.AdminPermissionValue) error {

	oRequest := &pbResourceModel.AdminPermissionAddOneInput{
		Value: pkgDomainToProto.AdminPermissionValue(oAdminPermission),
	}

	_, oErr := oSelf.ResourceModelClient.AdminPermission.AddOne(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminPermissionModel) ShowOneById(iId uint64) (*domain.AdminPermission, error) {

	oRequest := &pbResourceModel.AdminPermissionShowOneByIdInput{Id: uint64(iId)}
	oResponse, oErr := oSelf.ResourceModelClient.AdminPermission.ShowOneById(oSelf.Context, oRequest)

	if oErr != nil {
		return nil, oErr
	}

	oProtoAdminPermission := oResponse.GetAdminPermission()
	if oProtoAdminPermission.GetId() == 0 {
		return nil, nil
	}

	oAdminPermission := pkgProtoToDomain.AdminPermission(oProtoAdminPermission)

	return &oAdminPermission, nil
}

func (oSelf *AdminPermissionModel) EditOneById(oAdminPermission *domain.AdminPermissionValue, iId uint64) error {

	oRequest := &pbResourceModel.AdminPermissionEditOneByIdInput{
		Id:    uint64(iId),
		Value: pkgDomainToProto.AdminPermissionValue(oAdminPermission),
	}

	_, oErr := oSelf.ResourceModelClient.AdminPermission.EditOneById(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminPermissionModel) RemoveOneById(iId uint64) error {

	oRequest := &pbResourceModel.AdminPermissionRemoveOneByIdInput{Id: uint64(iId)}
	_, oErr := oSelf.ResourceModelClient.AdminPermission.RemoveOneById(oSelf.Context, oRequest)

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
		oAdminPermission := pkgProtoToDomain.AdminPermission(oOne)
		aAdminPermissions = append(aAdminPermissions, &oAdminPermission)
	}

	return aAdminPermissions, nil
}

func (oSelf *AdminPermissionModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.AdminPermissionTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminPermission.TotalByFilters(oSelf.Context, oRequest)

	iTotal := uint64(oResponse.GetTotal())

	return iTotal, oErr
}
