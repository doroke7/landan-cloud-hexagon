package inputApplicationResourceModel

import (
	"context"

	pb "example/pb"
	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgDomainToProto "example/pkg/domain_to_proto"
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

func (oSelf *AdminUserHandler) ShowOneByName(oContext context.Context, oReq *pbResourceModel.AdminUserShowOneByNameInput) (*pbResourceModel.AdminUserShowOneByNameOutput, error) {

	oAdminUser, err := oSelf.ModelAdminUserUsecase.ShowOneByName(oReq.Name)
	if err != nil {
		return nil, err
	}

	oProtoAdminUser := pkgDomainToProto.AdminUser(oAdminUser)

	return &pbResourceModel.AdminUserShowOneByNameOutput{
		AdminUser: oProtoAdminUser,
	}, nil

}

func (oSelf *AdminUserHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.AdminUserShowOneByIdInput) (*pbResourceModel.AdminUserShowOneByIdOutput, error) {

	oAdminUser, err := oSelf.ModelAdminUserUsecase.ShowOneById(uint64(oReq.Id))
	if err != nil {
		return nil, err
	}

	oProtoAdminUser := pkgDomainToProto.AdminUser(oAdminUser)

	return &pbResourceModel.AdminUserShowOneByIdOutput{
		AdminUser: oProtoAdminUser,
	}, nil

}

func (oSelf *AdminUserHandler) AddOne(oContext context.Context, oReq *pbResourceModel.AdminUserAddOneInput) (*pbResourceModel.AdminUserAddOneOutput, error) {

	oValue := oReq.GetVariable()

	oAdminUserValue := domain.AdminUserVariable{}
	if oValue != nil {
		oAdminUserValue.Name = oValue.Name
		oAdminUserValue.Password = oValue.Password
	}

	oErr := oSelf.ModelAdminUserUsecase.AddOne(&oAdminUserValue)

	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.AdminUserAddOneOutput{}, nil
}

func (oSelf *AdminUserHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.AdminUserEditOneByIdInput) (*pbResourceModel.AdminUserEditOneByIdOutput, error) {

	oValue := oReq.GetVariable()

	oAdminUserValue := domain.AdminUserVariable{}
	if oValue != nil {
		oAdminUserValue.Name = oValue.Name
		oAdminUserValue.Password = oValue.Password
	}

	oErr := oSelf.ModelAdminUserUsecase.EditOneById(&oAdminUserValue, uint64(oReq.Id))

	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.AdminUserEditOneByIdOutput{}, nil
}

func (oSelf *AdminUserHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.AdminUserRemoveOneByIdInput) (*pbResourceModel.AdminUserRemoveOneByIdOutput, error) {

	oErr := oSelf.ModelAdminUserUsecase.RemoveOneById(uint64(oReq.Id))

	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.AdminUserRemoveOneByIdOutput{}, nil
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
		return nil, oErr
	}

	aProtoAdminUsers := make([]*pb.AdminUser, 0, len(aAdminUsers))
	for _, oAdminUser := range aAdminUsers {
		oProtoAdminUser := pkgDomainToProto.AdminUser(oAdminUser)
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
		return nil, oErr
	}

	return &pbResourceModel.AdminUserTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
