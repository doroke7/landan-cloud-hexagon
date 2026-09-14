package inputApplicationResourceModel

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgDomainToProto "example/pkg/domain_to_proto"
	pkgInput "example/pkg/input"
)

type AdminRoleHandler struct {
	pbResourceModel.UnimplementedAdminRoleModelServer
	*inputApplicationResource.AbstractHandler
	ModelAdminRoleUsecase usecasePortAnyModel.AdminRoleUsecase
}

func NewAdminRoleHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oAdminRoleUsecase usecasePortAnyModel.AdminRoleUsecase) *AdminRoleHandler {
	return &AdminRoleHandler{
		AbstractHandler:       oAbstractHandler,
		ModelAdminRoleUsecase: oAdminRoleUsecase,
	}
}

func protoAdminRoleValueToDomainAdminRoleValue(oVariable *pbResource.AdminRoleValue) domain.AdminRoleVariable {
	var oValue domain.AdminRoleVariable
	if oVariable == nil {
		return oValue
	}

	oValue.Key = oVariable.Key
	oValue.Name = oVariable.Name

	return oValue
}

func (oSelf *AdminRoleHandler) AddOne(oContext context.Context, oReq *pbResourceModel.AdminRoleAddOneInput) (*pbResourceModel.AdminRoleAddOneOutput, error) {

	oAdminRoleValue := protoAdminRoleValueToDomainAdminRoleValue(oReq.GetValue())

	oErr := oSelf.ModelAdminRoleUsecase.AddOne(&oAdminRoleValue)

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AdminRoleAddOneOutput{}, nil
}

func (oSelf *AdminRoleHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.AdminRoleShowOneByIdInput) (*pbResourceModel.AdminRoleShowOneByIdOutput, error) {

	oAdminRole, oErr := oSelf.ModelAdminRoleUsecase.ShowOneById(uint(oReq.Id))
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	if oAdminRole == nil {
		return nil, nil
	}

	oProtoAdminRole := pkgDomainToProto.AdminRole(oAdminRole)

	return &pbResourceModel.AdminRoleShowOneByIdOutput{
		AdminRole: oProtoAdminRole,
	}, nil
}

func (oSelf *AdminRoleHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.AdminRoleEditOneByIdInput) (*pbResourceModel.AdminRoleEditOneByIdOutput, error) {

	oAdminRoleValue := protoAdminRoleValueToDomainAdminRoleValue(oReq.GetValue())

	oErr := oSelf.ModelAdminRoleUsecase.EditOneById(&oAdminRoleValue, uint64(oReq.Id))

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AdminRoleEditOneByIdOutput{}, nil
}

func (oSelf *AdminRoleHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.AdminRoleRemoveOneByIdInput) (*pbResourceModel.AdminRoleRemoveOneByIdOutput, error) {

	oErr := oSelf.ModelAdminRoleUsecase.RemoveOneById(uint64(oReq.Id))

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AdminRoleRemoveOneByIdOutput{}, nil
}

func (oSelf *AdminRoleHandler) ShowOnes(oContext context.Context, oReq *pbResourceModel.AdminRoleShowOnesInput) (*pbResourceModel.AdminRoleShowOnesOutput, error) {

	aAdminRoles, oErr := oSelf.ModelAdminRoleUsecase.ShowOnes()

	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	aProtoAdminRoles := make([]*pbResource.AdminRole, 0, len(aAdminRoles))
	for _, oAdminRole := range aAdminRoles {
		oProtoAdminRole := pkgDomainToProto.AdminRole(oAdminRole)
		aProtoAdminRoles = append(aProtoAdminRoles, oProtoAdminRole)
	}

	return &pbResourceModel.AdminRoleShowOnesOutput{
		AdminRoles: aProtoAdminRoles,
	}, nil
}

func (oSelf *AdminRoleHandler) ShowOnesByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceModel.AdminRoleShowOnesByFiltersWithSortersPaginationInput) (*pbResourceModel.AdminRoleShowOnesByFiltersWithSortersPaginationOutput, error) {

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

	aAdminRoles, oErr := oSelf.ModelAdminRoleUsecase.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	aProtoAdminRoles := make([]*pbResource.AdminRole, 0, len(aAdminRoles))
	for _, oAdminRole := range aAdminRoles {
		oProtoAdminRole := pkgDomainToProto.AdminRole(oAdminRole)
		aProtoAdminRoles = append(aProtoAdminRoles, oProtoAdminRole)
	}

	return &pbResourceModel.AdminRoleShowOnesByFiltersWithSortersPaginationOutput{
		AdminRoles: aProtoAdminRoles,
	}, nil
}

func (oSelf *AdminRoleHandler) TotalByFilters(oContext context.Context, oReq *pbResourceModel.AdminRoleTotalByFiltersInput) (*pbResourceModel.AdminRoleTotalByFiltersOutput, error) {

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

	iTotal, oErr := oSelf.ModelAdminRoleUsecase.TotalByFilters(aFilters)
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	return &pbResourceModel.AdminRoleTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
