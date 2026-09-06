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

type AdminUserHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceLogic.UnimplementedAdminUserLogicServer
	LogicAdminUserUsecase usecasePortAnyLogic.AdminUserUsecase
}

func NewAdminUserHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oAdminUserUsecase usecasePortAnyLogic.AdminUserUsecase) *AdminUserHandler {
	return &AdminUserHandler{
		AbstractHandler:       oAbstractHandler,
		LogicAdminUserUsecase: oAdminUserUsecase,
	}
}

func domainAdminRoleToProtoAdminRole(oAdminRole *domain.AdminRole) *pbResource.AdminRole {
	return &pbResource.AdminRole{
		Id:        uint64(oAdminRole.Id),
		Key:       oAdminRole.Key,
		Name:      oAdminRole.Name,
		CreatedAt: timestamppb.New(oAdminRole.CreatedAt),
		UpdatedAt: timestamppb.New(oAdminRole.UpdatedAt),
		DeletedAt: timestamppb.New(oAdminRole.DeletedAt),
	}
}

func domainAdminUserToProtoAdminUser(oAdminUser *domain.AdminUser) *pbResource.AdminUser {
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

func (oSelf *AdminUserHandler) ShowAdminUsersTotalByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceLogic.AdminUserShowAdminUsersTotalByFiltersWithSortersPaginationInput) (*pbResourceLogic.AdminUserShowAdminUsersTotalByFiltersWithSortersPaginationOutput, error) {

	iSize := uint(oReq.GetPagination().GetSize())
	iPage := uint(oReq.GetPagination().GetPage())
	oPagination := &pkgInput.Pagination{
		Size: &iSize,
		Page: &iPage,
	}

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

	aAdminUsers, iTotal, oErr := oSelf.LogicAdminUserUsecase.ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	aPbAdminUsers := make([]*pbResource.AdminUser, 0, len(aAdminUsers))
	for _, oAdminUser := range aAdminUsers {
		oProtoAdminUser := domainAdminUserToProtoAdminUser(oAdminUser)
		aPbAdminUsers = append(aPbAdminUsers, oProtoAdminUser)
	}

	return &pbResourceLogic.AdminUserShowAdminUsersTotalByFiltersWithSortersPaginationOutput{
		Total:      uint64(iTotal),
		AdminUsers: aPbAdminUsers,
	}, oErr
}

func (oSelf *AdminUserHandler) AddAminUser(oContext context.Context, oReq *pbResourceLogic.AdminUserAddAminUserInput) (*pbResourceLogic.AdminUserAddAminUserOutput, error) {

	aAdminRoleIds := make([]uint64, 0, len(oReq.GetValue().GetAdminRoleIds()))
	for _, iAdminRoleId := range oReq.GetValue().GetAdminRoleIds() {
		aAdminRoleIds = append(aAdminRoleIds, uint64(iAdminRoleId))
	}

	oValue := &domain.AdminUserValue{
		Name:         oReq.GetValue().Name,
		Password:     oReq.GetValue().Password,
		AdminRoleIds: aAdminRoleIds,
	}

	oErr := oSelf.LogicAdminUserUsecase.AddAminUser(oValue)

	return &pbResourceLogic.AdminUserAddAminUserOutput{}, oErr
}

func (oSelf *AdminUserHandler) EditAdminUserById(oContext context.Context, oReq *pbResourceLogic.AdminUserEditAdminUserByIdInput) (*pbResourceLogic.AdminUserEditAdminUserByIdOutput, error) {

	aAdminRoleIds := make([]uint64, 0, len(oReq.GetValue().GetAdminRoleIds()))
	for _, iAdminRoleId := range oReq.GetValue().GetAdminRoleIds() {
		aAdminRoleIds = append(aAdminRoleIds, uint64(iAdminRoleId))
	}

	oValue := &domain.AdminUserValue{
		Name:         oReq.GetValue().Name,
		Password:     oReq.GetValue().Password,
		AdminRoleIds: aAdminRoleIds,
	}

	oErr := oSelf.LogicAdminUserUsecase.EditAdminUserById(oValue, uint64(oReq.GetId()))

	return &pbResourceLogic.AdminUserEditAdminUserByIdOutput{}, oErr
}

func (oSelf *AdminUserHandler) RemoveAdminUserById(oContext context.Context, oReq *pbResourceLogic.AdminUserRemoveAdminUserByIdInput) (*pbResourceLogic.AdminUserRemoveAdminUserByIdOutput, error) {

	oErr := oSelf.LogicAdminUserUsecase.RemoveAdminUserById(oReq.GetId())

	return &pbResourceLogic.AdminUserRemoveAdminUserByIdOutput{}, oErr
}

func (oSelf *AdminUserHandler) ShowAdminUserById(oContext context.Context, oReq *pbResourceLogic.AdminUserShowAdminUserByIdInput) (*pbResourceLogic.AdminUserShowAdminUserByIdOutput, error) {

	oAdminUser, oErr := oSelf.LogicAdminUserUsecase.ShowAdminUserById(oReq.GetId())
	if oErr != nil {
		return nil, oErr
	}

	aPbAdminUsers := make([]*pbResource.AdminUser, 0, 1)
	if oAdminUser != nil {
		aPbAdminUsers = append(aPbAdminUsers, domainAdminUserToProtoAdminUser(oAdminUser))
	}

	oOutput := &pbResourceLogic.AdminUserShowAdminUserByIdOutput{
		AdminUsers: aPbAdminUsers,
		Total:      uint64(len(aPbAdminUsers)),
	}

	return oOutput, nil
}
