package inputApplicationResourceLogic

import (
	"context"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
)

type AdminRoleHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceLogic.UnimplementedAdminRoleLogicServer
	LogicAdminRoleUsecase usecasePortAnyLogic.AdminRoleUsecase
}

func NewAdminRoleHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oAdminRoleUsecase usecasePortAnyLogic.AdminRoleUsecase) *AdminRoleHandler {
	return &AdminRoleHandler{
		AbstractHandler:       oAbstractHandler,
		LogicAdminRoleUsecase: oAdminRoleUsecase,
	}
}

func (oSelf *AdminRoleHandler) AddAdminRole(oContext context.Context, oReq *pbResourceLogic.AdminRoleAddAdminRoleInput) (*pbResourceLogic.AdminRoleAddAdminRoleOutput, error) {

	oValue := &domain.AdminRoleVariable{
		Key:  oReq.GetVariable().Key,
		Name: oReq.GetVariable().Name,
	}

	oErr := oSelf.LogicAdminRoleUsecase.AddAdminRole(oValue)

	return &pbResourceLogic.AdminRoleAddAdminRoleOutput{}, oErr
}
