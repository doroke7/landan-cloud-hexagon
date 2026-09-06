package inputApplicationResourceLogic

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pbResource "example/pb/resource"
	pbResourceLogic "example/pb/resource/logic"
	pkgInput "example/pkg/input"
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

// ShowTree 從 usecase 拿到組好的 tree，攤平成平的 AdminPermissionGroup 回傳，
// client 端再自己組回 tree（組 tree 每個 adapter 各寫一份）。
func (oSelf *AdminPermissionGroupHandler) ShowTree(oContext context.Context, oReq *pbResourceLogic.AdminPermissionGroupShowTreeInput) (*pbResourceLogic.AdminPermissionGroupShowTreeOutput, error) {

	aRoots, oErr := oSelf.LogicAdminPermissionGroupUsecase.ShowTree()
	if oErr != nil {
		return nil, oErr
	}

	aNodes := make([]*pbResource.AdminPermissionGroup, 0, len(aRoots))

	var fnFlatten func(oNode domain.AdminPermissionGroup)
	fnFlatten = func(oNode domain.AdminPermissionGroup) {
		aNodes = append(aNodes, &pbResource.AdminPermissionGroup{
			Id:        uint64(oNode.Id),
			ParentId:  uint64(oNode.ParentId),
			Key:       oNode.Key,
			Name:      oNode.Name,
			CreatedAt: timestamppb.New(oNode.CreatedAt),
			UpdatedAt: timestamppb.New(oNode.UpdatedAt),
			DeletedAt: timestamppb.New(oNode.DeletedAt),
		})

		for _, oChild := range oNode.Children {
			fnFlatten(oChild)
		}
	}

	for _, oRoot := range aRoots {
		fnFlatten(*oRoot)
	}

	return &pbResourceLogic.AdminPermissionGroupShowTreeOutput{AdminPermissionGroups: aNodes}, nil
}

// domainAdminPermissionGroupToProtoAdminPermissionGroup 遞迴帶出 Parent / Children。
func domainAdminPermissionGroupToProtoAdminPermissionGroup(oNode *domain.AdminPermissionGroup) *pbResource.AdminPermissionGroup {
	if oNode == nil {
		return nil
	}

	oPb := &pbResource.AdminPermissionGroup{
		Id:        uint64(oNode.Id),
		ParentId:  uint64(oNode.ParentId),
		Key:       oNode.Key,
		Name:      oNode.Name,
		CreatedAt: timestamppb.New(oNode.CreatedAt),
		UpdatedAt: timestamppb.New(oNode.UpdatedAt),
		DeletedAt: timestamppb.New(oNode.DeletedAt),
		Parent:    domainAdminPermissionGroupToProtoAdminPermissionGroup(oNode.Parent),
	}

	for i := range oNode.Children {
		oPb.Children = append(oPb.Children, domainAdminPermissionGroupToProtoAdminPermissionGroup(&oNode.Children[i]))
	}

	return oPb
}

func (oSelf *AdminPermissionGroupHandler) ShowAdminPermissionGroupById(oContext context.Context, oReq *pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupByIdInput) (*pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupByIdOutput, error) {

	oAdminPermissionGroup, oErr := oSelf.LogicAdminPermissionGroupUsecase.ShowAdminPermissionGroupById(oReq.GetId())
	if oErr != nil {
		return nil, oErr
	}

	if oAdminPermissionGroup == nil {
		return nil, nil
	}

	return &pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupByIdOutput{
		AdminPermissionGroup: domainAdminPermissionGroupToProtoAdminPermissionGroup(oAdminPermissionGroup),
	}, nil
}

func (oSelf *AdminPermissionGroupHandler) ShowAdminPermissionGroups(oContext context.Context, oReq *pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupsInput) (*pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupsOutput, error) {

	aAdminPermissionGroups, oErr := oSelf.LogicAdminPermissionGroupUsecase.ShowAdminPermissionGroups()
	if oErr != nil {
		return nil, oErr
	}

	aNodes := make([]*pbResource.AdminPermissionGroup, 0, len(aAdminPermissionGroups))
	for _, oOne := range aAdminPermissionGroups {
		aNodes = append(aNodes, domainAdminPermissionGroupToProtoAdminPermissionGroup(oOne))
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
		aFilters = append(aFilters, &pkgInput.Filter{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oOne.GetValue().AsInterface(),
		})
	}

	aSorters := make([]*pkgInput.Sorter, 0, len(oReq.GetSorters()))
	for _, oOne := range oReq.GetSorters() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOrder := oOne.GetOrder()
		aSorters = append(aSorters, &pkgInput.Sorter{
			Field: &sField,
			Order: &sOrder,
		})
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
