package inputApplicationResourceModel

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceModel.UnimplementedAdminPermissionGroupModelServer
	ModelAdminPermissionGroupUsecase usecasePortAnyModel.AdminPermissionGroupUsecase
}

func NewAdminPermissionGroupHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oAdminPermissionGroupUsecase usecasePortAnyModel.AdminPermissionGroupUsecase) *AdminPermissionGroupHandler {
	return &AdminPermissionGroupHandler{
		AbstractHandler:                  oAbstractHandler,
		ModelAdminPermissionGroupUsecase: oAdminPermissionGroupUsecase,
	}
}

// domainAdminPermissionGroupToProtoAdminPermissionGroup 遞迴帶出 Parent / Children。
func domainAdminPermissionGroupToProtoAdminPermissionGroup(oAdminPermissionGroup *domain.AdminPermissionGroup) *pbResource.AdminPermissionGroup {
	if oAdminPermissionGroup == nil {
		return nil
	}

	oPb := &pbResource.AdminPermissionGroup{
		Id:        uint64(oAdminPermissionGroup.Id),
		ParentId:  uint64(oAdminPermissionGroup.ParentId),
		Key:       oAdminPermissionGroup.Key,
		Name:      oAdminPermissionGroup.Name,
		CreatedAt: timestamppb.New(oAdminPermissionGroup.CreatedAt),
		UpdatedAt: timestamppb.New(oAdminPermissionGroup.UpdatedAt),
		DeletedAt: timestamppb.New(oAdminPermissionGroup.DeletedAt),
		Parent:    domainAdminPermissionGroupToProtoAdminPermissionGroup(oAdminPermissionGroup.Parent),
	}

	for i := range oAdminPermissionGroup.Children {
		oPb.Children = append(oPb.Children, domainAdminPermissionGroupToProtoAdminPermissionGroup(&oAdminPermissionGroup.Children[i]))
	}

	return oPb
}

func protoAdminPermissionGroupVariableToDomainAdminPermissionGroup(oVariable *pbResourceModel.AdminPermissionGroupVariable) domain.AdminPermissionGroup {
	var oAdminPermissionGroup domain.AdminPermissionGroup
	if oVariable == nil {
		return oAdminPermissionGroup
	}

	oAdminPermissionGroup.ParentId = oVariable.GetParentId()
	oAdminPermissionGroup.Key = oVariable.GetKey()
	oAdminPermissionGroup.Name = oVariable.GetName()

	return oAdminPermissionGroup
}

func (oSelf *AdminPermissionGroupHandler) AddOne(oContext context.Context, oReq *pbResourceModel.AdminPermissionGroupAddOneInput) (*pbResourceModel.AdminPermissionGroupAddOneOutput, error) {

	oAdminPermissionGroup := protoAdminPermissionGroupVariableToDomainAdminPermissionGroup(oReq.GetVariable())

	oErr := oSelf.ModelAdminPermissionGroupUsecase.AddOne(&oAdminPermissionGroup)

	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.AdminPermissionGroupAddOneOutput{
		AdminPermissionGroup: &pbResource.AdminPermissionGroup{
			ParentId: oAdminPermissionGroup.ParentId,
			Key:      oAdminPermissionGroup.Key,
			Name:     oAdminPermissionGroup.Name,
		},
	}, nil
}

func (oSelf *AdminPermissionGroupHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.AdminPermissionGroupShowOneByIdInput) (*pbResourceModel.AdminPermissionGroupShowOneByIdOutput, error) {

	oAdminPermissionGroup, oErr := oSelf.ModelAdminPermissionGroupUsecase.ShowOneById(oReq.GetId())

	if oErr != nil {
		return nil, oErr
	}

	if oAdminPermissionGroup == nil {
		return nil, nil
	}

	return &pbResourceModel.AdminPermissionGroupShowOneByIdOutput{
		AdminPermissionGroup: domainAdminPermissionGroupToProtoAdminPermissionGroup(oAdminPermissionGroup),
	}, nil
}

func (oSelf *AdminPermissionGroupHandler) ShowOnesByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceModel.AdminPermissionGroupShowOnesByFiltersWithSortersPaginationInput) (*pbResourceModel.AdminPermissionGroupShowOnesByFiltersWithSortersPaginationOutput, error) {

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

	aAdminPermissionGroups, oErr := oSelf.ModelAdminPermissionGroupUsecase.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	if oErr != nil {
		return nil, oErr
	}

	aPbAdminPermissionGroups := make([]*pbResource.AdminPermissionGroup, 0, len(aAdminPermissionGroups))
	for _, oAdminPermissionGroup := range aAdminPermissionGroups {
		aPbAdminPermissionGroups = append(aPbAdminPermissionGroups, domainAdminPermissionGroupToProtoAdminPermissionGroup(oAdminPermissionGroup))
	}

	return &pbResourceModel.AdminPermissionGroupShowOnesByFiltersWithSortersPaginationOutput{
		AdminPermissionGroups: aPbAdminPermissionGroups,
	}, nil
}

func (oSelf *AdminPermissionGroupHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.AdminPermissionGroupEditOneByIdInput) (*pbResourceModel.AdminPermissionGroupEditOneByIdOutput, error) {

	oAdminPermissionGroup := protoAdminPermissionGroupVariableToDomainAdminPermissionGroup(oReq.GetVariable())

	oErr := oSelf.ModelAdminPermissionGroupUsecase.EditOneById(&oAdminPermissionGroup, oReq.GetId())
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.AdminPermissionGroupEditOneByIdOutput{
		Status: true,
	}, nil
}

func (oSelf *AdminPermissionGroupHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.AdminPermissionGroupRemoveOneByIdInput) (*pbResourceModel.AdminPermissionGroupRemoveOneByIdOutput, error) {

	oErr := oSelf.ModelAdminPermissionGroupUsecase.RemoveOneById(oReq.GetId())
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.AdminPermissionGroupRemoveOneByIdOutput{
		Status: true,
	}, nil
}

func (oSelf *AdminPermissionGroupHandler) TotalByFilters(oContext context.Context, oReq *pbResourceModel.AdminPermissionGroupTotalByFiltersInput) (*pbResourceModel.AdminPermissionGroupTotalByFiltersOutput, error) {

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

	iTotal, oErr := oSelf.ModelAdminPermissionGroupUsecase.TotalByFilters(aFilters)

	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.AdminPermissionGroupTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
