package inputApplicationResourceEvent

import (
	"context"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyEvent "example/internal/usecase/port/any/event"
	pbResourceEvent "example/pb/resource/event"
)

type AdminUserHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceEvent.UnimplementedAdminUserEventServer
	usecasePortAnyEvent.AdminUserUsecase
}

func NewAdminUserHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oAdminUserUsecase usecasePortAnyEvent.AdminUserUsecase) *AdminUserHandler {
	return &AdminUserHandler{
		AbstractHandler:  oAbstractHandler,
		AdminUserUsecase: oAdminUserUsecase,
	}
}

func (oSelf *AdminUserHandler) AddOne(oContext context.Context, oReq *pbResourceEvent.AdminUserEventAddOneInput) (*pbResourceEvent.AdminUserEventAddOneOutput, error) {

	oErr := oSelf.AdminUserUsecase.AddOne(&domain.AdminUser{
		Id:       uint(oReq.GetId()),
		Name:     oReq.GetName(),
		Password: oReq.GetPassword(),
	})
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceEvent.AdminUserEventAddOneOutput{}, nil
}
