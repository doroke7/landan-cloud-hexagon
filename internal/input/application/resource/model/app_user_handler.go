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

type AppUserHandler struct {
	pbResourceModel.UnimplementedAppUserModelServer
	*inputApplicationResource.AbstractHandler
	ModelAppUserUsecase usecasePortAnyModel.AppUserUsecase
}

func NewAppUserHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oAppUserUsecase usecasePortAnyModel.AppUserUsecase) *AppUserHandler {
	return &AppUserHandler{
		AbstractHandler:     oAbstractHandler,
		ModelAppUserUsecase: oAppUserUsecase,
	}
}

func (oSelf *AppUserHandler) AddAppUser(oContext context.Context, oReq *pbResourceModel.AppUserAddOneInput) (*pbResourceModel.AppUserAddOneOutput, error) {

	oVariable := oReq.GetVariable()

	oAppUserValue := domain.AppUserValue{}
	if oVariable != nil {
		oAppUserValue.Name = oVariable.Name
		oAppUserValue.Password = oVariable.Password
	}

	if _, oErr := oSelf.ModelAppUserUsecase.AddOne(&oAppUserValue); oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	var sName string
	if oAppUserValue.Name != nil {
		sName = *oAppUserValue.Name
	}

	return &pbResourceModel.AppUserAddOneOutput{
		Name: sName,
	}, nil
}

func (oSelf *AppUserHandler) ShowOneByName(oContext context.Context, oReq *pbResourceModel.AppUserShowOneByNameInput) (*pbResourceModel.AppUserShowOneByNameOutput, error) {

	oAppUser, oErr := oSelf.ModelAppUserUsecase.ShowOneByName(oReq.Name)
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	return &pbResourceModel.AppUserShowOneByNameOutput{
		AppUser: &pbResource.AppUser{
			Id:       uint32(oAppUser.Id),
			Name:     oAppUser.Name,
			Password: oAppUser.Password,
			Balance:  uint32(oAppUser.Balance),
		},
	}, nil
}

func (oSelf *AppUserHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.AppUserShowOneByIdInput) (*pbResourceModel.AppUserShowOneByIdOutput, error) {

	oAppUser, oErr := oSelf.ModelAppUserUsecase.ShowOneById(uint(oReq.Id))
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	return &pbResourceModel.AppUserShowOneByIdOutput{
		AppUser: &pbResource.AppUser{
			Id:       uint32(oAppUser.Id),
			Name:     oAppUser.Name,
			Password: oAppUser.Password,
			Balance:  uint32(oAppUser.Balance),
		},
	}, nil
}

func (oSelf *AppUserHandler) IncreaseBalance(oContext context.Context, oReq *pbResourceModel.AppUserIncreaseBalanceInput) (*pbResourceModel.AppUserIncreaseBalanceOutput, error) {

	bResult, oErr := oSelf.ModelAppUserUsecase.IncreaseBalance(uint(oReq.Id), uint(oReq.Amount))
	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AppUserIncreaseBalanceOutput{
		Status: bResult,
	}, nil
}
