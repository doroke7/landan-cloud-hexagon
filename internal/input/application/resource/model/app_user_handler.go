package inputApplicationResourceModel

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
)

type AppUserHandler struct {
	pbResourceModel.UnimplementedAppUserModelServer
	*inputApplicationResource.AbstractHandler
	usecasePortAnyModel.AppUserUsecase
}

func NewAppUserHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oAppUserUsecase usecasePortAnyModel.AppUserUsecase) *AppUserHandler {
	return &AppUserHandler{
		AbstractHandler: oAbstractHandler,
		AppUserUsecase:  oAppUserUsecase,
	}
}

func (oSelf *AppUserHandler) AddAppUser(oContext context.Context, oReq *pbResourceModel.AppUserAddOneInput) (*pbResourceModel.AppUserAddOneOutput, error) {

	oAppUserValue := domain.AppUserValue{
		Name: &oReq.Name,
	}

	if _, oErr := oSelf.AppUserUsecase.AddOne(&oAppUserValue); oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AppUserAddOneOutput{
		Name: oReq.Name,
	}, nil
}

func (oSelf *AppUserHandler) ShowOneByName(oContext context.Context, oReq *pbResourceModel.AppUserShowOneByNameInput) (*pbResourceModel.AppUserShowOneByNameOutput, error) {

	oAppUser, oErr := oSelf.AppUserUsecase.ShowOneByName(oReq.Name)
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	return &pbResourceModel.AppUserShowOneByNameOutput{
		Id:       uint32(oAppUser.Id),
		Name:     oAppUser.Name,
		Password: oAppUser.Password,
		Balance:  uint32(oAppUser.Balance),
	}, nil
}

func (oSelf *AppUserHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.AppUserShowOneByIdInput) (*pbResourceModel.AppUserShowOneByIdOutput, error) {

	oAppUser, oErr := oSelf.AppUserUsecase.ShowOneById(uint(oReq.Id))
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	return &pbResourceModel.AppUserShowOneByIdOutput{
		Id:       uint32(oAppUser.Id),
		Name:     oAppUser.Name,
		Password: oAppUser.Password,
		Balance:  uint32(oAppUser.Balance),
	}, nil
}

func (oSelf *AppUserHandler) IncreaseBalance(oContext context.Context, oReq *pbResourceModel.AppUserIncreaseBalanceInput) (*pbResourceModel.AppUserIncreaseBalanceOutput, error) {

	bResult, oErr := oSelf.AppUserUsecase.IncreaseBalance(uint(oReq.Id), uint(oReq.Amount))
	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.AppUserIncreaseBalanceOutput{
		Status: bResult,
	}, nil
}
