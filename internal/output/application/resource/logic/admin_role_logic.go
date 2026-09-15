package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
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
	return nil
}
