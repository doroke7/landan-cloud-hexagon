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
	var aAdminPermissionIds []uint64
	if oVariable.AdminPermissionIds != nil {
		aAdminPermissionIds = *oVariable.AdminPermissionIds
	}

	oRequest := &pbResourceLogic.AdminRoleAddAdminRoleInput{
		Variable: &pb.AdminRoleVariable{
			Key:                oVariable.Key,
			Name:               oVariable.Name,
			AdminPermissionIds: aAdminPermissionIds,
		},
	}

	_, oErr := oSelf.ResourceLogicClient.AdminRole.AddAdminRole(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminRoleLogic) EditAdminRoleById(oVariable *domain.AdminRoleVariable, iId uint64) error {
	return nil
}
