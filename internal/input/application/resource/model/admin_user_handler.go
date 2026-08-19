package input_application_resource_model

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
)

type AdminUserHandler struct {
	pbResourceModel.UnimplementedAdminUserModelServer
	*inputApplicationResource.AbstractHandler
	usecasePortAnyModel.AdminUserUsecase
}

func NewAdminUserHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oAdminUserUsecase usecasePortAnyModel.AdminUserUsecase) *AdminUserHandler {
	return &AdminUserHandler{
		AbstractHandler:  oAbstractHandler,
		AdminUserUsecase: oAdminUserUsecase,
	}
}

func (oSelf *AdminUserHandler) ShowOneByName(oContext context.Context, oReq *pbResourceModel.AdminUserShowOneByNameInput) (*pbResourceModel.AdminUserShowOneByNameOutput, error) {

	oAdminUser, err := oSelf.AdminUserUsecase.ShowOneByName(oReq.Name)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pbResourceModel.AdminUserShowOneByNameOutput{
		Id:       int32(oAdminUser.Id),
		Name:     oAdminUser.Name,
		Password: oAdminUser.Password,
	}, nil

}

func (oSelf *AdminUserHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.AdminUserShowOneByIdInput) (*pbResourceModel.AdminUserShowOneByIdOutput, error) {

	oAdminUser, err := oSelf.AdminUserUsecase.ShowOneById(uint(oReq.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pbResourceModel.AdminUserShowOneByIdOutput{
		Id:       uint32(oAdminUser.Id),
		Name:     oAdminUser.Name,
		Password: oAdminUser.Password,
	}, nil

}

func (oSelf *AdminUserHandler) AddOne(oContext context.Context, oReq *pbResourceModel.AdminUserAddOneInput) (*pbResourceModel.AdminUserAddOneOutput, error) {

	oAdminUserValue := domain.AdminUserValue{
		Name:     oReq.Name,
		Password: oReq.Password,
	}

	bResult, oErr := oSelf.AdminUserUsecase.AddOne(&oAdminUserValue)

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AdminUserAddOneOutput{
		Status: bResult,
	}, nil
}
