package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResource "example/pb"
	pbResourceLogic "example/pb/resource/logic"
	pkgDomainToProto "example/pkg/domain_to_proto"
	pkgInput "example/pkg/input"
	pkgProtoToDomain "example/pkg/proto_to_domain"
)

type AdminPermissionGroupLogic struct {
	*AbstractLogic
}

func NewAdminPermissionGroupLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminPermissionGroupLogic {
	oLogic := &AdminPermissionGroupLogic{
		AbstractLogic: oAbstractLogic,
	}

	return oLogic
}

func (oSelf *AdminPermissionGroupLogic) AddAdminPermissionGroup(oValue *domain.AdminPermissionGroupVariable) error {

	aAdminPermissions := make([]*pbResource.AdminPermissionVariable, 0, len(oValue.AdminPermissions))
	for _, oAdminPermissionValue := range oValue.AdminPermissions {
		oProtoAdminPermissionValue := pkgDomainToProto.AdminPermissionVariable(oAdminPermissionValue)
		if oProtoAdminPermissionValue == nil {
			continue
		}

		aAdminPermissions = append(aAdminPermissions, oProtoAdminPermissionValue)
	}

	oProtoValue := &pbResource.AdminPermissionGroupVariable{
		Key:              oValue.Key,
		Name:             oValue.Name,
		AdminPermissions: aAdminPermissions,
	}

	oRequest := &pbResourceLogic.AdminPermissionGroupAddAdminPermissionGroupInput{
		Variable: oProtoValue,
	}

	_, oErr := oSelf.ResourceLogicClient.AdminPermissionGroup.AddAdminPermissionGroup(oSelf.Context, oRequest)

	return oErr
}
func (oSelf *AdminPermissionGroupLogic) EditAdminPermissionGroupById(oValue *domain.AdminPermissionGroupVariable, iId uint64) error {

	aAdminPermissions := make([]*pbResource.AdminPermissionVariable, 0, len(oValue.AdminPermissions))
	for _, oAdminPermissionValue := range oValue.AdminPermissions {
		oProtoAdminPermissionValue := pkgDomainToProto.AdminPermissionVariable(oAdminPermissionValue)
		if oProtoAdminPermissionValue == nil {
			continue
		}

		aAdminPermissions = append(aAdminPermissions, oProtoAdminPermissionValue)
	}

	oProtoValue := &pbResource.AdminPermissionGroupVariable{
		Key:              oValue.Key,
		Name:             oValue.Name,
		AdminPermissions: aAdminPermissions,
	}

	oRequest := &pbResourceLogic.AdminPermissionGroupEditAdminPermissionGroupByIdInput{
		Variable: oProtoValue,
		Id:       iId,
	}

	_, oErr := oSelf.ResourceLogicClient.AdminPermissionGroup.EditAdminPermissionGroupById(oSelf.Context, oRequest)

	return oErr
}

// ShowTree gRPC 回來就是巢狀好的 tree（每個節點帶 Children），直接遞迴轉成 domain。
func (oSelf *AdminPermissionGroupLogic) ShowTree() ([]*domain.AdminPermissionGroup, error) {

	oResponse, oErr := oSelf.ResourceLogicClient.AdminPermissionGroup.ShowTree(oSelf.Context, &pbResourceLogic.AdminPermissionGroupShowTreeInput{})
	if oErr != nil {
		return nil, oErr
	}

	aRoots := make([]*domain.AdminPermissionGroup, 0, len(oResponse.GetAdminPermissionGroups()))
	for _, oOne := range oResponse.GetAdminPermissionGroups() {
		oRoot := pkgProtoToDomain.AdminPermissionGroup(oOne)
		aRoots = append(aRoots, &oRoot)
	}

	return aRoots, nil
}

func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroupById(iId uint64) (*domain.AdminPermissionGroup, error) {

	oRequest := &pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupByIdInput{Id: iId}
	oResponse, oErr := oSelf.ResourceLogicClient.AdminPermissionGroup.ShowAdminPermissionGroupById(oSelf.Context, oRequest)
	if oErr != nil {
		return nil, oErr
	}

	oProto := oResponse.GetAdminPermissionGroup()
	if oProto.GetId() == 0 {
		return nil, nil
	}

	oAdminPermissionGroup := pkgProtoToDomain.AdminPermissionGroup(oProto)

	return &oAdminPermissionGroup, nil
}

func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroups() ([]*domain.AdminPermissionGroup, error) {

	oResponse, oErr := oSelf.ResourceLogicClient.AdminPermissionGroup.ShowAdminPermissionGroups(
		oSelf.Context,
		&pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupsInput{},
	)
	if oErr != nil {
		return nil, oErr
	}

	aAdminPermissionGroups := make([]*domain.AdminPermissionGroup, 0, len(oResponse.GetAdminPermissionGroups()))
	for _, oOne := range oResponse.GetAdminPermissionGroups() {
		oAdminPermissionGroup := pkgProtoToDomain.AdminPermissionGroup(oOne)
		aAdminPermissionGroups = append(aAdminPermissionGroups, &oAdminPermissionGroup)
	}

	return aAdminPermissionGroups, nil
}

func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, uint64, error) {

	oRequest := &pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupsTotalByFiltersWithSortersPaginationInput{
		Filters: oSelf.ToFilters(aFilters),
		Sorters: oSelf.ToSorters(aSorters),
	}

	if oPagination != nil {
		oRequest.Pagination = oSelf.ToPagination(oPagination)
	}

	oResponse, oErr := oSelf.ResourceLogicClient.AdminPermissionGroup.ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(oSelf.Context, oRequest)
	if oErr != nil {
		return nil, 0, oErr
	}

	aAdminPermissionGroups := make([]*domain.AdminPermissionGroup, 0, len(oResponse.GetAdminPermissionGroups()))
	for _, oOne := range oResponse.GetAdminPermissionGroups() {
		oAdminPermissionGroup := pkgProtoToDomain.AdminPermissionGroup(oOne)
		aAdminPermissionGroups = append(aAdminPermissionGroups, &oAdminPermissionGroup)
	}

	iTotal := uint64(oResponse.GetTotal())
	return aAdminPermissionGroups, iTotal, nil
}
