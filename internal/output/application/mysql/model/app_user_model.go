package mysql

import (
	"gorm.io/gorm"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AppUserModel struct {
	*AbstractModel
}

func NewAppUserModel(oAbstractModel *AbstractModel) outputPortAnyModel.AppUserModel {
	return &AppUserModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *AppUserModel) IncreaseBalance(id uint, amount uint) (bool, error) {
	var oAppUserRow domain.AppUserRow

	if err := oSelf.DB.WithContext(oSelf.Context).Model(&oAppUserRow).
		Where("id = ?", id).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
		return false, err
	}

	return true, nil
}
