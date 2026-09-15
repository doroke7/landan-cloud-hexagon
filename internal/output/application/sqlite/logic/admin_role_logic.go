package outputApplicationSqliteLogic

import (
	"errors"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgUtility "example/pkg/utility"
)

type AdminRoleLogic struct {
	*AbstractLogic
}

func NewAdminRoleLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminRoleLogic {
	return &AdminRoleLogic{
		AbstractLogic: oAbstractLogic,
	}
}

func (oSelf *AdminRoleLogic) AddAdminRole(oVariable *domain.AdminRoleVariable) error {
	oColumns, oErr := pkgUtility.StructToMap(oVariable)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminRole{}).
		Create(oColumns)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		oZeroRowsError := errors.New("0 rows inserted")

		return oZeroRowsError
	}

	return nil
}
