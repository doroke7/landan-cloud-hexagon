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

func domainAdminUserToProtoAdminUser(oAdminUser *domain.AdminUser) *pbResource.AdminUser {
	return &pbResource.AdminUser{
		Id:        uint32(oAdminUser.Id),
		Name:      oAdminUser.Name,
		Password:  oAdminUser.Password,
		CreatedAt: timestamppb.New(oAdminUser.CreatedAt),
		UpdatedAt: timestamppb.New(oAdminUser.UpdatedAt),
		DeletedAt: timestamppb.New(oAdminUser.DeletedAt),
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
		oProtoAdminUser.Password = "" // 列表不外洩密碼
		aPbAdminUsers = append(aPbAdminUsers, oProtoAdminUser)
	}

	return &pbResourceLogic.AdminUserShowAdminUsersTotalByFiltersWithSortersPaginationOutput{
		Total:      iTotal,
		AdminUsers: aPbAdminUsers,
	}, oErr
}
