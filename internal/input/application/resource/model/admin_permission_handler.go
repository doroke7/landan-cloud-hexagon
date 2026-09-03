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
)

type AdminPermissionHandler struct {
	pbResourceModel.UnimplementedAdminPermissionModelServer
	*inputApplicationResource.AbstractHandler
	usecasePortAnyModel.AdminPermissionUsecase
}

func NewAdminPermissionHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oAdminPermissionUsecase usecasePortAnyModel.AdminPermissionUsecase) *AdminPermissionHandler {
	return &AdminPermissionHandler{
		AbstractHandler:        oAbstractHandler,
		AdminPermissionUsecase: oAdminPermissionUsecase,
	}
}

func domainAdminPermissionToProtoAdminPermission(oAdminPermission *domain.AdminPermission) *pbResource.AdminPermission {
	if oAdminPermission == nil {
		return nil
	}

	return &pbResource.AdminPermission{
		Id:        uint32(oAdminPermission.Id),
		Type:      uint32(oAdminPermission.Type),
		Key:       oAdminPermission.Key,
		Name:      oAdminPermission.Name,
		CreatedAt: timestamppb.New(oAdminPermission.CreatedAt),
		UpdatedAt: timestamppb.New(oAdminPermission.UpdatedAt),
		DeletedAt: timestamppb.New(oAdminPermission.DeletedAt),
	}
}

func protoAdminPermissionVariableToDomainAdminPermissionValue(oVariable *pbResourceModel.AdminPermissionVariable) domain.AdminPermissionValue {
	var oValue domain.AdminPermissionValue
	if oVariable == nil {
		return oValue
	}

	oValue.Key = oVariable.Key
	oValue.Name = oVariable.Name

	if oVariable.Type != nil {
		iType := uint8(*oVariable.Type)
		oValue.Type = &iType
	}

	return oValue
}

func (oSelf *AdminPermissionHandler) AddOne(oContext context.Context, oReq *pbResourceModel.AdminPermissionAddOneInput) (*pbResourceModel.AdminPermissionAddOneOutput, error) {

	oAdminPermissionValue := protoAdminPermissionVariableToDomainAdminPermissionValue(oReq.GetVariable())

	bResult, oErr := oSelf.AdminPermissionUsecase.AddOne(&oAdminPermissionValue)

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AdminPermissionAddOneOutput{
		Status: bResult,
	}, nil
}

func (oSelf *AdminPermissionHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.AdminPermissionShowOneByIdInput) (*pbResourceModel.AdminPermissionShowOneByIdOutput, error) {

	oAdminPermission, oErr := oSelf.AdminPermissionUsecase.ShowOneById(uint(oReq.Id))
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	if oAdminPermission == nil {
		return nil, nil
	}

	return &pbResourceModel.AdminPermissionShowOneByIdOutput{
		AdminPermission: domainAdminPermissionToProtoAdminPermission(oAdminPermission),
	}, nil
}

func (oSelf *AdminPermissionHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.AdminPermissionEditOneByIdInput) (*pbResourceModel.AdminPermissionEditOneByIdOutput, error) {

	oAdminPermissionValue := protoAdminPermissionVariableToDomainAdminPermissionValue(oReq.GetVariable())

	bResult, oErr := oSelf.AdminPermissionUsecase.EditOneById(&oAdminPermissionValue, uint(oReq.Id))

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AdminPermissionEditOneByIdOutput{
		Status: bResult,
	}, nil
}

func (oSelf *AdminPermissionHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.AdminPermissionRemoveOneByIdInput) (*pbResourceModel.AdminPermissionRemoveOneByIdOutput, error) {

	bResult, oErr := oSelf.AdminPermissionUsecase.RemoveOneById(uint(oReq.Id))

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AdminPermissionRemoveOneByIdOutput{
		Status: bResult,
	}, nil
}

func (oSelf *AdminPermissionHandler) ShowOnesByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceModel.AdminPermissionShowOnesByFiltersWithSortersPaginationInput) (*pbResourceModel.AdminPermissionShowOnesByFiltersWithSortersPaginationOutput, error) {

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

	aAdminPermissions, oErr := oSelf.AdminPermissionUsecase.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	aProtoAdminPermissions := make([]*pbResource.AdminPermission, 0, len(aAdminPermissions))
	for _, oAdminPermission := range aAdminPermissions {
		aProtoAdminPermissions = append(aProtoAdminPermissions, domainAdminPermissionToProtoAdminPermission(oAdminPermission))
	}

	return &pbResourceModel.AdminPermissionShowOnesByFiltersWithSortersPaginationOutput{
		AdminPermissions: aProtoAdminPermissions,
	}, nil
}

func (oSelf *AdminPermissionHandler) TotalByFilters(oContext context.Context, oReq *pbResourceModel.AdminPermissionTotalByFiltersInput) (*pbResourceModel.AdminPermissionTotalByFiltersOutput, error) {

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

	iTotal, oErr := oSelf.AdminPermissionUsecase.TotalByFilters(aFilters)
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	return &pbResourceModel.AdminPermissionTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
