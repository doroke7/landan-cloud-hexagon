package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResource "example/pb/resource"
	pbResourceLogic "example/pb/resource/logic"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupLogic struct {
	*AbstractLogic
}

func NewAdminPermissionGroupLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminPermissionGroupLogic {
	return &AdminPermissionGroupLogic{
		AbstractLogic: oAbstractLogic,
	}
}

func (oSelf *AdminPermissionGroupLogic) AddAdminPermissionGroup(oValue *domain.AdminPermissionGroupValue) error {

	aAdminPermissions := make([]*pbResourceLogic.AdminPermissionValue, 0, len(oValue.AdminPermissions))
	for _, oAdminPermissionValue := range oValue.AdminPermissions {
		oProtoAdminPermissionValue := domainAdminPermissionValueToProtoAdminPermissionValue(oAdminPermissionValue)
		if oProtoAdminPermissionValue == nil {
			continue
		}

		aAdminPermissions = append(aAdminPermissions, oProtoAdminPermissionValue)
	}

	oProtoValue := &pbResourceLogic.AdminPermissionGroupValue{
		Key:              oValue.Key,
		Name:             oValue.Name,
		AdminPermissions: aAdminPermissions,
	}

	oRequest := &pbResourceLogic.AdminPermissionGroupAddAdminPermissionGroupInput{
		Value: oProtoValue,
	}

	_, oErr := oSelf.ResourceLogicClient.AdminPermissionGroup.AddAdminPermissionGroup(oSelf.Context, oRequest)

	return oErr
}

func domainAdminPermissionValueToProtoAdminPermissionValue(oValue *domain.AdminPermissionValue) *pbResourceLogic.AdminPermissionValue {
	if oValue == nil {
		return nil
	}

	var oType *uint64
	if oValue.Type != nil {
		iType := uint64(*oValue.Type)
		oType = &iType
	}

	oProtoAdminPermissionValue := &pbResourceLogic.AdminPermissionValue{
		Id:   oValue.Id,
		Type: oType,
		Key:  oValue.Key,
		Name: oValue.Name,
	}

	return oProtoAdminPermissionValue
}

// ShowTree gRPC 回來就是巢狀好的 tree（每個節點帶 Children），直接遞迴轉成 domain。
func (oSelf *AdminPermissionGroupLogic) ShowTree() ([]*domain.AdminPermissionGroup, error) {
	oResponse, oErr := oSelf.ResourceLogicClient.AdminPermissionGroup.ShowTree(oSelf.Context, &pbResourceLogic.AdminPermissionGroupShowTreeInput{})
	if oErr != nil {
		return nil, oErr
	}

	aRoots := make([]*domain.AdminPermissionGroup, 0, len(oResponse.GetAdminPermissionGroups()))
	for _, oOne := range oResponse.GetAdminPermissionGroups() {
		oRoot := protoAdminPermissionGroupToDomainAdminPermissionGroup(oOne)
		aRoots = append(aRoots, &oRoot)
	}

	return aRoots, nil
}

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
		// 沒有子節點時也回 []（非 nil），避免 JSON 出現 children: null
		Children: make([]domain.AdminPermissionGroup, 0, len(oProto.GetChildren())),
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

func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroupById(iId uint64) (*domain.AdminPermissionGroup, error) {

	oResponse, oErr := oSelf.ResourceLogicClient.AdminPermissionGroup.ShowAdminPermissionGroupById(
		oSelf.Context,
		&pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupByIdInput{Id: iId},
	)
	if oErr != nil {
		return nil, oErr
	}

	oProto := oResponse.GetAdminPermissionGroup()
	if oProto.GetId() == 0 {
		return nil, nil
	}

	oAdminPermissionGroup := protoAdminPermissionGroupToDomainAdminPermissionGroup(oProto)

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
		oAdminPermissionGroup := protoAdminPermissionGroupToDomainAdminPermissionGroup(oOne)
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
		oAdminPermissionGroup := protoAdminPermissionGroupToDomainAdminPermissionGroup(oOne)
		aAdminPermissionGroups = append(aAdminPermissionGroups, &oAdminPermissionGroup)
	}

	iTotal := uint64(oResponse.GetTotal())
	return aAdminPermissionGroups, iTotal, nil
}
