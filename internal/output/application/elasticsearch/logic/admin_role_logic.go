package outputApplicationElasticsearchLogic

import (
	"strconv"
	"time"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
)

type AdminRoleLogic struct {
	*AbstractLogic
	Index string
}

func NewAdminRoleLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminRoleLogic {
	return &AdminRoleLogic{
		AbstractLogic: oAbstractLogic,
		Index:         oAbstractLogic.IndexName("admin_roles"),
	}
}

func (oSelf *AdminRoleLogic) AddAdminRole(oVariable *domain.AdminRoleVariable) error {
	iId, oErr := oSelf.NextId("admin_role")
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oDoc := &domain.AdminRole{
		Id:        uint64(iId),
		CreatedAt: oNow,
		UpdatedAt: oNow,
	}

	if oVariable.Key != nil {
		oDoc.Key = *oVariable.Key
	}
	if oVariable.Name != nil {
		oDoc.Name = *oVariable.Name
	}

	if oErr := oSelf.IndexOne(oSelf.Index, strconv.FormatUint(uint64(iId), 10), oDoc); oErr != nil {
		return oErr
	}

	return nil
}
