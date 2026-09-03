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
		AdminUser: &pbResource.AdminUser{
			Id:       uint32(oAdminUser.Id),
			Name:     oAdminUser.Name,
			Password: oAdminUser.Password,
		},
	}, nil

}

func (oSelf *AdminUserHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.AdminUserShowOneByIdInput) (*pbResourceModel.AdminUserShowOneByIdOutput, error) {

	oAdminUser, err := oSelf.AdminUserUsecase.ShowOneById(uint(oReq.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pbResourceModel.AdminUserShowOneByIdOutput{
		AdminUser: &pbResource.AdminUser{
			Id:       uint32(oAdminUser.Id),
			Name:     oAdminUser.Name,
			Password: oAdminUser.Password,
		},
	}, nil

}

func (oSelf *AdminUserHandler) AddOne(oContext context.Context, oReq *pbResourceModel.AdminUserAddOneInput) (*pbResourceModel.AdminUserAddOneOutput, error) {

	oVariable := oReq.GetVariable()

	oAdminUserValue := domain.AdminUserValue{}
	if oVariable != nil {
		oAdminUserValue.Name = oVariable.Name
		oAdminUserValue.Password = oVariable.Password
	}

	bResult, oErr := oSelf.AdminUserUsecase.AddOne(&oAdminUserValue)

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AdminUserAddOneOutput{
		Status: bResult,
	}, nil
}

func (oSelf *AdminUserHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.AdminUserEditOneByIdInput) (*pbResourceModel.AdminUserEditOneByIdOutput, error) {

	oVariable := oReq.GetVariable()

	oAdminUserValue := domain.AdminUserValue{}
	if oVariable != nil {
		oAdminUserValue.Name = oVariable.Name
		oAdminUserValue.Password = oVariable.Password
	}

	bResult, oErr := oSelf.AdminUserUsecase.EditOneById(&oAdminUserValue, uint(oReq.Id))

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AdminUserEditOneByIdOutput{
		Status: bResult,
	}, nil
}
