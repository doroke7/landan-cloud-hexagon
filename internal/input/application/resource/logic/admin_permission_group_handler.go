package inputApplicationResourceLogic

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pbResource "example/pb/resource"
	pbResourceLogic "example/pb/resource/logic"
	pkgDomainToProto "example/pkg/domain_to_proto"
	pkgInput "example/pkg/input"
	pkgProtoToDomain "example/pkg/proto_to_domain"
)

type AdminPermissionGroupHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceLogic.UnimplementedAdminPermissionGroupLogicServer
	LogicAdminPermissionGroupUsecase usecasePortAnyLogic.AdminPermissionGroupUsecase
}

func NewAdminPermissionGroupHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oAdminPermissionGroupUsecase usecasePortAnyLogic.AdminPermissionGroupUsecase) *AdminPermissionGroupHandler {
	return &AdminPermissionGroupHandler{
		AbstractHandler:                  oAbstractHandler,
		LogicAdminPermissionGroupUsecase: oAdminPermissionGroupUsecase,
	}
}

func (oSelf *AdminPermissionGroupHandler) AddAdminPermissionGroup(oContext context.Context, oReq *pbResourceLogic.AdminPermissionGroupAddAdminPermissionGroupInput) (*pbResourceLogic.AdminPermissionGroupAddAdminPermissionGroupOutput, error) {

	var oValue domain.AdminPermissionGroupVariable
	if oVariable := oReq.GetVariable(); oVariable != nil {
		oValue.Key = oVariable.Key
		oValue.Name = oVariable.Name

		aAdminPermissions := make([]*domain.AdminPermissionVariable, 0, len(oVariable.GetAdminPermissions()))
		for _, oProtoAdminPermissionValue := range oVariable.GetAdminPermissions() {
			if oProtoAdminPermissionValue == nil {
				continue
			}

			oAdminPermissionValue := pkgProtoToDomain.AdminPermissionVariable(oProtoAdminPermissionValue)
			aAdminPermissions = append(aAdminPermissions, &oAdminPermissionValue)
		}

		oValue.AdminPermissions = aAdminPermissions
	}

	oErr := oSelf.LogicAdminPermissionGroupUsecase.AddAdminPermissionGroup(&oValue)
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceLogic.AdminPermissionGroupAddAdminPermissionGroupOutput{
		Status: true,
	}, nil
}

func (oSelf *AdminPermissionGroupHandler) EditAdminPermissionGroupById(oContext context.Context, oReq *pbResourceLogic.AdminPermissionGroupEditAdminPermissionGroupByIdInput) (*pbResourceLogic.AdminPermissionGroupEditAdminPermissionGroupByIdOutput, error) {

	var oValue domain.AdminPermissionGroupVariable
	if oVariable := oReq.GetVariable(); oVariable != nil {
		oValue.Key = oVariable.Key
		oValue.Name = oVariable.Name

		aAdminPermissions := make([]*domain.AdminPermissionVariable, 0, len(oVariable.GetAdminPermissions()))
		for _, oProtoAdminPermissionValue := range oVariable.GetAdminPermissions() {
			if oProtoAdminPermissionValue == nil {
				continue
			}

			oAdminPermissionValue := pkgProtoToDomain.AdminPermissionVariable(oProtoAdminPermissionValue)
			aAdminPermissions = append(aAdminPermissions, &oAdminPermissionValue)
		}

		oValue.AdminPermissions = aAdminPermissions
	}

	iId := oReq.GetId()
	oErr := oSelf.LogicAdminPermissionGroupUsecase.EditAdminPermissionGroupById(&oValue, iId)
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceLogic.AdminPermissionGroupEditAdminPermissionGroupByIdOutput{}, nil
}

// ShowTree 從 usecase 拿到組好的 tree，攤平成平的 AdminPermissionGroup 回傳，
// client 端再自己組回 tree（組 tree 每個 adapter 各寫一份）。
func (oSelf *AdminPermissionGroupHandler) ShowTree(oContext context.Context, oReq *pbResourceLogic.AdminPermissionGroupShowTreeInput) (*pbResourceLogic.AdminPermissionGroupShowTreeOutput, error) {

	aTrees, oErr := oSelf.LogicAdminPermissionGroupUsecase.ShowTree()
	if oErr != nil {
		return nil, oErr
	}

	aProtoAdminPermissionGroups := make([]*pbResource.AdminPermissionGroup, 0, len(aTrees))
	for _, oAdminPermissionGroup := range aTrees {
		oProtoAdminPermissionGroup := pkgDomainToProto.AdminPermissionGroup(oAdminPermissionGroup)
		aProtoAdminPermissionGroups = append(aProtoAdminPermissionGroups, oProtoAdminPermissionGroup)
	}

	return &pbResourceLogic.AdminPermissionGroupShowTreeOutput{AdminPermissionGroups: aProtoAdminPermissionGroups}, nil
}

func (oSelf *AdminPermissionGroupHandler) ShowAdminPermissionGroupById(oContext context.Context, oReq *pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupByIdInput) (*pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupByIdOutput, error) {

	iId := oReq.GetId()
	oAdminPermissionGroup, oErr := oSelf.LogicAdminPermissionGroupUsecase.ShowAdminPermissionGroupById(iId)
	if oErr != nil {
		return nil, oErr
	}

	if oAdminPermissionGroup == nil {
		return nil, nil
	}

	oProtoAdminPermissionGroup := pkgDomainToProto.AdminPermissionGroup(oAdminPermissionGroup)

	return &pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupByIdOutput{
		AdminPermissionGroup: oProtoAdminPermissionGroup,
	}, nil
}

func (oSelf *AdminPermissionGroupHandler) ShowAdminPermissionGroups(oContext context.Context, oReq *pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupsInput) (*pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupsOutput, error) {

	aAdminPermissionGroups, oErr := oSelf.LogicAdminPermissionGroupUsecase.ShowAdminPermissionGroups()
	if oErr != nil {
		return nil, oErr
	}

	aNodes := make([]*pbResource.AdminPermissionGroup, 0, len(aAdminPermissionGroups))
	for _, oOne := range aAdminPermissionGroups {
		oNode := pkgDomainToProto.AdminPermissionGroup(oOne)
		aNodes = append(aNodes, oNode)
	}

	return &pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupsOutput{AdminPermissionGroups: aNodes}, nil
}

func (oSelf *AdminPermissionGroupHandler) ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupsTotalByFiltersWithSortersPaginationInput) (*pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupsTotalByFiltersWithSortersPaginationOutput, error) {

	aFilters := make([]*pkgInput.Filter, 0, len(oReq.GetFilters()))
	for _, oOne := range oReq.GetFilters() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOperator := oOne.GetOperator()
		oValue := oOne.GetValue().AsInterface()
		oFilter := &pkgInput.Filter{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oValue,
		}
		aFilters = append(aFilters, oFilter)
	}

	aSorters := make([]*pkgInput.Sorter, 0, len(oReq.GetSorters()))
	for _, oOne := range oReq.GetSorters() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOrder := oOne.GetOrder()
		oSorter := &pkgInput.Sorter{
			Field: &sField,
			Order: &sOrder,
		}
		aSorters = append(aSorters, oSorter)
	}

	iSize := uint(oReq.GetPagination().GetSize())
	iPage := uint(oReq.GetPagination().GetPage())
	oPagination := &pkgInput.Pagination{
		Size: &iSize,
		Page: &iPage,
	}

	aAdminPermissionGroups, iTotal, oErr := oSelf.LogicAdminPermissionGroupUsecase.ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, oErr
	}

	aNodes := make([]*pbResource.AdminPermissionGroup, 0, len(aAdminPermissionGroups))
	for _, oOne := range aAdminPermissionGroups {
		aNodes = append(aNodes, &pbResource.AdminPermissionGroup{
			Id:        uint64(oOne.Id),
			ParentId:  uint64(oOne.ParentId),
			Key:       oOne.Key,
			Name:      oOne.Name,
			CreatedAt: timestamppb.New(oOne.CreatedAt),
			UpdatedAt: timestamppb.New(oOne.UpdatedAt),
			DeletedAt: timestamppb.New(oOne.DeletedAt),
		})
	}

	return &pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupsTotalByFiltersWithSortersPaginationOutput{
		AdminPermissionGroups: aNodes,
		Total:                 iTotal,
	}, nil
}
