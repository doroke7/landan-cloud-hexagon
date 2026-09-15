package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pb "example/pb"
	pbResourceLogic "example/pb/resource/logic"
)

type AdminRoleLogic struct {
	*AbstractLogic
}

func NewAdminRoleLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminRoleLogic {
	oLogic := &AdminRoleLogic{
		AbstractLogic: oAbstractLogic,
	}

	return oLogic
}

func (oSelf *AdminRoleLogic) AddAdminRole(oVariable *domain.AdminRoleVariable) error {
	oRequest := &pbResourceLogic.AdminRoleAddAdminRoleInput{
		Variable: &pb.AdminRoleVariable{
			Key:  oVariable.Key,
			Name: oVariable.Name,
		},
	}

	_, oErr := oSelf.ResourceLogicClient.AdminRole.AddAdminRole(oSelf.Context, oRequest)

	return oErr
}
