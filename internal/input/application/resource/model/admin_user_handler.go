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

type AdminUserHandler struct {
	pbResourceModel.UnimplementedAdminUserModelServer
	*inputApplicationResource.AbstractHandler
	ModelAdminUserUsecase usecasePortAnyModel.AdminUserUsecase
}

func NewAdminUserHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oAdminUserUsecase usecasePortAnyModel.AdminUserUsecase) *AdminUserHandler {
	return &AdminUserHandler{
		AbstractHandler:       oAbstractHandler,
		ModelAdminUserUsecase: oAdminUserUsecase,
	}
}

func domainAdminUserToProtoAdminUser(oAdminUser *domain.AdminUser) *pbResource.AdminUser {
	if oAdminUser == nil {
		return nil
	}

	aAdminRoles := make([]*pbResource.AdminRole, 0, len(oAdminUser.AdminRoles))
	for i := range oAdminUser.AdminRoles {
		aAdminRoles = append(aAdminRoles, domainAdminRoleToProtoAdminRole(&oAdminUser.AdminRoles[i]))
	}

	return &pbResource.AdminUser{
		Id:         uint64(oAdminUser.Id),
		Name:       oAdminUser.Name,
		Password:   oAdminUser.Password,
		CreatedAt:  timestamppb.New(oAdminUser.CreatedAt),
		UpdatedAt:  timestamppb.New(oAdminUser.UpdatedAt),
		DeletedAt:  timestamppb.New(oAdminUser.DeletedAt),
		AdminRoles: aAdminRoles,
	}
}

func (oSelf *AdminUserHandler) ShowOneByName(oContext context.Context, oReq *pbResourceModel.AdminUserShowOneByNameInput) (*pbResourceModel.AdminUserShowOneByNameOutput, error) {

	oAdminUser, err := oSelf.ModelAdminUserUsecase.ShowOneByName(oReq.Name)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pbResourceModel.AdminUserShowOneByNameOutput{
		AdminUser: domainAdminUserToProtoAdminUser(oAdminUser),
	}, nil

}

func (oSelf *AdminUserHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.AdminUserShowOneByIdInput) (*pbResourceModel.AdminUserShowOneByIdOutput, error) {

	oAdminUser, err := oSelf.ModelAdminUserUsecase.ShowOneById(uint64(oReq.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pbResourceModel.AdminUserShowOneByIdOutput{
		AdminUser: domainAdminUserToProtoAdminUser(oAdminUser),
	}, nil

}

func (oSelf *AdminUserHandler) AddOne(oContext context.Context, oReq *pbResourceModel.AdminUserAddOneInput) (*pbResourceModel.AdminUserAddOneOutput, error) {

	oVariable := oReq.GetVariable()

	oAdminUserValue := domain.AdminUserValue{}
	if oVariable != nil {
		oAdminUserValue.Name = oVariable.Name
		oAdminUserValue.Password = oVariable.Password
	}

	oErr := oSelf.ModelAdminUserUsecase.AddOne(&oAdminUserValue)

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AdminUserAddOneOutput{
		Status: true,
	}, nil
}

func (oSelf *AdminUserHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.AdminUserEditOneByIdInput) (*pbResourceModel.AdminUserEditOneByIdOutput, error) {

	oVariable := oReq.GetVariable()

	oAdminUserValue := domain.AdminUserValue{}
	if oVariable != nil {
		oAdminUserValue.Name = oVariable.Name
		oAdminUserValue.Password = oVariable.Password
	}

	oErr := oSelf.ModelAdminUserUsecase.EditOneById(&oAdminUserValue, uint64(oReq.Id))

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AdminUserEditOneByIdOutput{
		Status: true,
	}, nil
}

func (oSelf *AdminUserHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.AdminUserRemoveOneByIdInput) (*pbResourceModel.AdminUserRemoveOneByIdOutput, error) {

	oErr := oSelf.ModelAdminUserUsecase.RemoveOneById(uint64(oReq.Id))

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AdminUserRemoveOneByIdOutput{
		Status: true,
	}, nil
}

func (oSelf *AdminUserHandler) ShowOnesByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceModel.AdminUserShowOnesByFiltersWithSortersPaginationInput) (*pbResourceModel.AdminUserShowOnesByFiltersWithSortersPaginationOutput, error) {

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

	aAdminUsers, oErr := oSelf.ModelAdminUserUsecase.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	aProtoAdminUsers := make([]*pbResource.AdminUser, 0, len(aAdminUsers))
	for _, oAdminUser := range aAdminUsers {
		oProtoAdminUser := domainAdminUserToProtoAdminUser(oAdminUser)
		aProtoAdminUsers = append(aProtoAdminUsers, oProtoAdminUser)
	}

	return &pbResourceModel.AdminUserShowOnesByFiltersWithSortersPaginationOutput{
		AdminUsers: aProtoAdminUsers,
	}, nil
}

func (oSelf *AdminUserHandler) TotalByFilters(oContext context.Context, oReq *pbResourceModel.AdminUserTotalByFiltersInput) (*pbResourceModel.AdminUserTotalByFiltersOutput, error) {

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

	iTotal, oErr := oSelf.ModelAdminUserUsecase.TotalByFilters(aFilters)
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	return &pbResourceModel.AdminUserTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
