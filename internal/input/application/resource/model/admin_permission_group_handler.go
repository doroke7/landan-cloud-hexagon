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

func protoAdminPermissionGroupValueToDomainAdminPermissionGroupValue(oVariable *pbResourceModel.AdminPermissionGroupValue) domain.AdminPermissionGroupValue {
	var oValue domain.AdminPermissionGroupValue
	if oVariable == nil {
		return oValue
	}

	oValue.Key = oVariable.Key
	oValue.Name = oVariable.Name

	return oValue
}

func (oSelf *AdminPermissionGroupHandler) AddOne(oContext context.Context, oReq *pbResourceModel.AdminPermissionGroupAddOneInput) (*pbResourceModel.AdminPermissionGroupAddOneOutput, error) {

	oValue := protoAdminPermissionGroupValueToDomainAdminPermissionGroupValue(oReq.GetValue())

	oErr := oSelf.ModelAdminPermissionGroupUsecase.AddOne(&oValue)

	if oErr != nil {
		return nil, oErr
	}

	oOut := &pbResource.AdminPermissionGroup{}
	if oValue.Key != nil {
		oOut.Key = *oValue.Key
	}
	if oValue.Name != nil {
		oOut.Name = *oValue.Name
	}

	return &pbResourceModel.AdminPermissionGroupAddOneOutput{
		AdminPermissionGroup: oOut,
	}, nil
}

func (oSelf *AdminPermissionGroupHandler) ShowOnesByParentId(oContext context.Context, oReq *pbResourceModel.AdminPermissionGroupShowOnesByParentIdInput) (*pbResourceModel.AdminPermissionGroupShowOnesByParentIdOutput, error) {

	aAdminPermissionGroups, oErr := oSelf.ModelAdminPermissionGroupUsecase.ShowOnesByParentId(oReq.GetParentId())

	if oErr != nil {
		return nil, oErr
	}

	aPbAdminPermissionGroups := make([]*pbResource.AdminPermissionGroup, 0, len(aAdminPermissionGroups))
	for _, oAdminPermissionGroup := range aAdminPermissionGroups {
		aPbAdminPermissionGroups = append(aPbAdminPermissionGroups, domainAdminPermissionGroupToProtoAdminPermissionGroup(oAdminPermissionGroup))
	}

	return &pbResourceModel.AdminPermissionGroupShowOnesByParentIdOutput{
		AdminPermissionGroups: aPbAdminPermissionGroups,
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

	oValue := protoAdminPermissionGroupValueToDomainAdminPermissionGroupValue(oReq.GetValue())

	oErr := oSelf.ModelAdminPermissionGroupUsecase.EditOneById(&oValue, oReq.GetId())
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
