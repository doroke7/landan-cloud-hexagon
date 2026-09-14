package inputApplicationResourceModel

import (
	"context"

	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"

	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgDomainToProto "example/pkg/domain_to_proto"
	pkgInput "example/pkg/input"
	pkgProtoToDomain "example/pkg/proto_to_domain"
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

func (oSelf *AdminPermissionGroupHandler) AddOne(oContext context.Context, oReq *pbResourceModel.AdminPermissionGroupAddOneInput) (*pbResourceModel.AdminPermissionGroupAddOneOutput, error) {

	oReqValue := oReq.GetVariable()
	oValue := pkgProtoToDomain.AdminPermissionGroupVariable(oReqValue)

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

	iParentId := oReq.GetParentId()
	aAdminPermissionGroups, oErr := oSelf.ModelAdminPermissionGroupUsecase.ShowOnesByParentId(iParentId)

	if oErr != nil {
		return nil, oErr
	}

	aPbAdminPermissionGroups := make([]*pbResource.AdminPermissionGroup, 0, len(aAdminPermissionGroups))
	for _, oAdminPermissionGroup := range aAdminPermissionGroups {
		oPbAdminPermissionGroup := pkgDomainToProto.AdminPermissionGroup(oAdminPermissionGroup)
		aPbAdminPermissionGroups = append(aPbAdminPermissionGroups, oPbAdminPermissionGroup)
	}

	return &pbResourceModel.AdminPermissionGroupShowOnesByParentIdOutput{
		AdminPermissionGroups: aPbAdminPermissionGroups,
	}, nil
}

func (oSelf *AdminPermissionGroupHandler) ShowOnesByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceModel.AdminPermissionGroupShowOnesByFiltersWithSortersPaginationInput) (*pbResourceModel.AdminPermissionGroupShowOnesByFiltersWithSortersPaginationOutput, error) {

	aFiltersInput := oReq.GetFilters()
	aFilters := make([]*pkgInput.Filter, 0, len(aFiltersInput))
	for _, oOne := range aFiltersInput {
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

	aSortersInput := oReq.GetSorters()
	aSorters := make([]*pkgInput.Sorter, 0, len(aSortersInput))
	for _, oOne := range aSortersInput {
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

	oPaginationInput := oReq.GetPagination()
	iSize := uint(oPaginationInput.GetSize())
	iPage := uint(oPaginationInput.GetPage())
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
		oPbAdminPermissionGroup := pkgDomainToProto.AdminPermissionGroup(oAdminPermissionGroup)
		aPbAdminPermissionGroups = append(aPbAdminPermissionGroups, oPbAdminPermissionGroup)
	}

	return &pbResourceModel.AdminPermissionGroupShowOnesByFiltersWithSortersPaginationOutput{
		AdminPermissionGroups: aPbAdminPermissionGroups,
	}, nil
}

func (oSelf *AdminPermissionGroupHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.AdminPermissionGroupEditOneByIdInput) (*pbResourceModel.AdminPermissionGroupEditOneByIdOutput, error) {

	oReqValue := oReq.GetVariable()
	oValue := pkgProtoToDomain.AdminPermissionGroupVariable(oReqValue)

	iId := oReq.GetId()
	oErr := oSelf.ModelAdminPermissionGroupUsecase.EditOneById(&oValue, iId)
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.AdminPermissionGroupEditOneByIdOutput{}, nil
}

func (oSelf *AdminPermissionGroupHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.AdminPermissionGroupRemoveOneByIdInput) (*pbResourceModel.AdminPermissionGroupRemoveOneByIdOutput, error) {

	iId := oReq.GetId()
	oErr := oSelf.ModelAdminPermissionGroupUsecase.RemoveOneById(iId)
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.AdminPermissionGroupRemoveOneByIdOutput{}, nil
}

func (oSelf *AdminPermissionGroupHandler) TotalByFilters(oContext context.Context, oReq *pbResourceModel.AdminPermissionGroupTotalByFiltersInput) (*pbResourceModel.AdminPermissionGroupTotalByFiltersOutput, error) {

	aFiltersInput := oReq.GetFilters()
	aFilters := make([]*pkgInput.Filter, 0, len(aFiltersInput))
	for _, oOne := range aFiltersInput {
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

	iTotal, oErr := oSelf.ModelAdminPermissionGroupUsecase.TotalByFilters(aFilters)

	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.AdminPermissionGroupTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
