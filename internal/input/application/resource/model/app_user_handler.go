package inputApplicationResourceModel

import (
	"context"

	pb "example/pb"
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

	oValue := oReq.GetVariable()

	oAppUserValue := domain.AppUserVariable{}
	if oValue != nil {
		oAppUserValue.Name = oValue.Name
		oAppUserValue.Password = oValue.Password
	}

	oErr := oSelf.ModelAppUserUsecase.AddOne(&oAppUserValue)

	if oErr != nil {
		return nil, oErr
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
		return nil, oErr
	}

	return &pbResourceModel.AppUserShowOneByNameOutput{
		AppUser: &pb.AppUser{
			Id:       uint64(oAppUser.Id),
			Name:     oAppUser.Name,
			Password: oAppUser.Password,
			Balance:  uint64(oAppUser.Balance),
		},
	}, nil
}

func (oSelf *AppUserHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.AppUserShowOneByIdInput) (*pbResourceModel.AppUserShowOneByIdOutput, error) {

	oAppUser, oErr := oSelf.ModelAppUserUsecase.ShowOneById(uint(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.AppUserShowOneByIdOutput{
		AppUser: &pb.AppUser{
			Id:       uint64(oAppUser.Id),
			Name:     oAppUser.Name,
			Password: oAppUser.Password,
			Balance:  uint64(oAppUser.Balance),
		},
	}, nil
}

func (oSelf *AppUserHandler) IncreaseBalance(oContext context.Context, oReq *pbResourceModel.AppUserIncreaseBalanceInput) (*pbResourceModel.AppUserIncreaseBalanceOutput, error) {

	oErr := oSelf.ModelAppUserUsecase.IncreaseBalance(uint64(oReq.Id), uint64(oReq.Amount))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.AppUserIncreaseBalanceOutput{}, nil
}
