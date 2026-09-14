package inputApplicationResourceModel

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
	pkgProtoToDomain "example/pkg/proto_to_domain"
)

type AdminPermissionHandler struct {
	pbResourceModel.UnimplementedAdminPermissionModelServer
	*inputApplicationResource.AbstractHandler
	ModelAdminPermissionUsecase usecasePortAnyModel.AdminPermissionUsecase
}

func NewAdminPermissionHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oAdminPermissionUsecase usecasePortAnyModel.AdminPermissionUsecase) *AdminPermissionHandler {
	return &AdminPermissionHandler{
		AbstractHandler:             oAbstractHandler,
		ModelAdminPermissionUsecase: oAdminPermissionUsecase,
	}
}

func domainAdminPermissionToProtoAdminPermission(oAdminPermission *domain.AdminPermission) *pbResource.AdminPermission {
	if oAdminPermission == nil {
		return nil
	}

	return &pbResource.AdminPermission{
		Id:                     uint64(oAdminPermission.Id),
		AdminPermissionGroupId: oAdminPermission.AdminPermissionGroupId,
		Type:                   uint64(oAdminPermission.Type),
		Key:                    oAdminPermission.Key,
		Name:                   oAdminPermission.Name,
		CreatedAt:              timestamppb.New(oAdminPermission.CreatedAt),
		UpdatedAt:              timestamppb.New(oAdminPermission.UpdatedAt),
		DeletedAt:              timestamppb.New(oAdminPermission.DeletedAt),
	}
}

func (oSelf *AdminPermissionHandler) AddOne(oContext context.Context, oReq *pbResourceModel.AdminPermissionAddOneInput) (*pbResourceModel.AdminPermissionAddOneOutput, error) {

	oReqValue := oReq.GetVariable()
	oAdminPermissionValue := pkgProtoToDomain.AdminPermissionVariable(oReqValue)

	oErr := oSelf.ModelAdminPermissionUsecase.AddOne(&oAdminPermissionValue)

	if oErr != nil {
		sError := oErr.Error()
		oStatusError := status.Error(codes.Aborted, sError)

		return nil, oStatusError
	}

	return &pbResourceModel.AdminPermissionAddOneOutput{
		Status: true,
	}, nil
}

func (oSelf *AdminPermissionHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.AdminPermissionShowOneByIdInput) (*pbResourceModel.AdminPermissionShowOneByIdOutput, error) {

	oAdminPermission, oErr := oSelf.ModelAdminPermissionUsecase.ShowOneById(uint64(oReq.Id))
	if oErr != nil {
		sError := oErr.Error()
		oStatusError := status.Error(codes.NotFound, sError)

		return nil, oStatusError
	}

	if oAdminPermission == nil {
		return nil, nil
	}

	oProtoAdminPermission := domainAdminPermissionToProtoAdminPermission(oAdminPermission)

	return &pbResourceModel.AdminPermissionShowOneByIdOutput{
		AdminPermission: oProtoAdminPermission,
	}, nil
}

func (oSelf *AdminPermissionHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.AdminPermissionEditOneByIdInput) (*pbResourceModel.AdminPermissionEditOneByIdOutput, error) {

	oReqValue := oReq.GetVariable()
	oAdminPermissionValue := pkgProtoToDomain.AdminPermissionVariable(oReqValue)

	oErr := oSelf.ModelAdminPermissionUsecase.EditOneById(&oAdminPermissionValue, uint64(oReq.Id))

	if oErr != nil {
		sError := oErr.Error()
		oStatusError := status.Error(codes.Aborted, sError)

		return nil, oStatusError
	}

	return &pbResourceModel.AdminPermissionEditOneByIdOutput{}, nil
}

func (oSelf *AdminPermissionHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.AdminPermissionRemoveOneByIdInput) (*pbResourceModel.AdminPermissionRemoveOneByIdOutput, error) {

	oErr := oSelf.ModelAdminPermissionUsecase.RemoveOneById(uint64(oReq.Id))

	if oErr != nil {
		sError := oErr.Error()
		oStatusError := status.Error(codes.Aborted, sError)

		return nil, oStatusError
	}

	return &pbResourceModel.AdminPermissionRemoveOneByIdOutput{}, nil
}

func (oSelf *AdminPermissionHandler) ShowOnesByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceModel.AdminPermissionShowOnesByFiltersWithSortersPaginationInput) (*pbResourceModel.AdminPermissionShowOnesByFiltersWithSortersPaginationOutput, error) {

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

	aAdminPermissions, oErr := oSelf.ModelAdminPermissionUsecase.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	if oErr != nil {
		sError := oErr.Error()
		oStatusError := status.Error(codes.NotFound, sError)

		return nil, oStatusError
	}

	aProtoAdminPermissions := make([]*pbResource.AdminPermission, 0, len(aAdminPermissions))
	for _, oAdminPermission := range aAdminPermissions {
		oProtoAdminPermission := domainAdminPermissionToProtoAdminPermission(oAdminPermission)
		aProtoAdminPermissions = append(aProtoAdminPermissions, oProtoAdminPermission)
	}

	return &pbResourceModel.AdminPermissionShowOnesByFiltersWithSortersPaginationOutput{
		AdminPermissions: aProtoAdminPermissions,
	}, nil
}

func (oSelf *AdminPermissionHandler) TotalByFilters(oContext context.Context, oReq *pbResourceModel.AdminPermissionTotalByFiltersInput) (*pbResourceModel.AdminPermissionTotalByFiltersOutput, error) {

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

	iTotal, oErr := oSelf.ModelAdminPermissionUsecase.TotalByFilters(aFilters)
	if oErr != nil {
		sError := oErr.Error()
		oStatusError := status.Error(codes.NotFound, sError)

		return nil, oStatusError
	}

	return &pbResourceModel.AdminPermissionTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
