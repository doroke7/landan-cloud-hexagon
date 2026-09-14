package inputApplicationResourceLogic

import (
	"context"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pb "example/pb"
	pbResourceLogic "example/pb/resource/logic"
	pkgDomainToProto "example/pkg/domain_to_proto"
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

	aPbAdminUsers := make([]*pb.AdminUser, 0, len(aAdminUsers))
	for _, oAdminUser := range aAdminUsers {
		oProtoAdminUser := pkgDomainToProto.AdminUser(oAdminUser)
		aPbAdminUsers = append(aPbAdminUsers, oProtoAdminUser)
	}

	return &pbResourceLogic.AdminUserShowAdminUsersTotalByFiltersWithSortersPaginationOutput{
		Total:      uint64(iTotal),
		AdminUsers: aPbAdminUsers,
	}, oErr
}

func (oSelf *AdminUserHandler) AddAminUser(oContext context.Context, oReq *pbResourceLogic.AdminUserAddAminUserInput) (*pbResourceLogic.AdminUserAddAminUserOutput, error) {

	aAdminRoleIds := make([]uint64, 0, len(oReq.GetVariable().GetAdminRoleIds()))
	for _, iAdminRoleId := range oReq.GetVariable().GetAdminRoleIds() {
		aAdminRoleIds = append(aAdminRoleIds, uint64(iAdminRoleId))
	}

	oValue := &domain.AdminUserVariable{
		Name:         oReq.GetVariable().Name,
		Password:     oReq.GetVariable().Password,
		AdminRoleIds: aAdminRoleIds,
	}

	oErr := oSelf.LogicAdminUserUsecase.AddAminUser(oValue)

	return &pbResourceLogic.AdminUserAddAminUserOutput{}, oErr
}

func (oSelf *AdminUserHandler) EditAdminUserById(oContext context.Context, oReq *pbResourceLogic.AdminUserEditAdminUserByIdInput) (*pbResourceLogic.AdminUserEditAdminUserByIdOutput, error) {

	aAdminRoleIds := make([]uint64, 0, len(oReq.GetVariable().GetAdminRoleIds()))
	for _, iAdminRoleId := range oReq.GetVariable().GetAdminRoleIds() {
		aAdminRoleIds = append(aAdminRoleIds, uint64(iAdminRoleId))
	}

	oValue := &domain.AdminUserVariable{
		Name:         oReq.GetVariable().Name,
		Password:     oReq.GetVariable().Password,
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

	aPbAdminUsers := make([]*pb.AdminUser, 0, 1)
	if oAdminUser != nil {
		oProtoAdminUser := pkgDomainToProto.AdminUser(oAdminUser)
		aPbAdminUsers = append(aPbAdminUsers, oProtoAdminUser)
	}

	oOutput := &pbResourceLogic.AdminUserShowAdminUserByIdOutput{
		AdminUsers: aPbAdminUsers,
		Total:      uint64(len(aPbAdminUsers)),
	}

	return oOutput, nil
}
