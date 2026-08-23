package mysql

import (
	"errors"

	"gorm.io/gorm"

	domain "example/internal/domain"
	mysqlBase "example/internal/output/application/mysql"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkg "example/pkg"
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

func (oSelf *AppUserModel) ShowOneByName(sName string) (*domain.AppUser, error) {
	var oAppUserRow domain.AppUserRow

	if oErr := oSelf.DB.WithContext(oSelf.Context).Where("name = ?", sName).First(&oAppUserRow).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, errors.New("資料不存在")
		}
		return nil, oErr
	}

	oAppUser := domain.AppUserRowToAppUser(&oAppUserRow)

	return oAppUser, nil
}

func (oSelf *AppUserModel) ShowOneById(iId uint) (*domain.AppUser, error) {
	var oAppUserRow domain.AppUserRow

	if oErr := oSelf.DB.WithContext(oSelf.Context).First(&oAppUserRow, iId).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, errors.New("資料不存在")
		}
		return nil, oErr
	}

	oAppUser := domain.AppUserRowToAppUser(&oAppUserRow)

	return oAppUser, nil
}

func (oSelf *AppUserModel) AddOne(oValue *domain.AppUserValue) (bool, error) {
	var oAppUserRow domain.AppUserRow

	oColumns, oErr := pkg.StructToMap(oValue)
	if oErr != nil {
		return false, oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).Model(&oAppUserRow).Create(oColumns)
	if oResult.Error != nil {
		return false, oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return false, errors.New("新增0筆")
	}

	return true, nil
}
