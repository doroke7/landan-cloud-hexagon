package mysql

import (
	"gorm.io/gorm"

	domain "example/internal/domain"
	mysqlBase "example/internal/output/application/mysql"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AppUserModel struct {
	*mysqlBase.AbstractMysql
}

func NewAppUserModel(oAbstractModel *mysqlBase.AbstractMysql) outputPortAnyModel.AppUserModel {
	return &AppUserModel{
		AbstractMysql: oAbstractModel,
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
