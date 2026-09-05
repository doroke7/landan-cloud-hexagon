package outputApplicationMysqlModel

import (
	"errors"

	"gorm.io/gorm"

	domain "example/internal/domain"
	mysqlBase "example/internal/output/application/mysql"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgUtility "example/pkg/utility"
)

type AppUserModel struct {
	*mysqlBase.AbstractMysql
}

func NewAppUserModel(oAbstractModel *mysqlBase.AbstractMysql) outputPortAnyModel.AppUserModel {
	return &AppUserModel{
		AbstractMysql: oAbstractModel,
	}
}

func (oSelf *AppUserModel) IncreaseBalance(id uint, amount uint64) error {
	if err := oSelf.DB.WithContext(oSelf.Context).Model(&domain.AppUser{}).
		Where("id = ?", id).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
		return err
	}

	return nil
}

func (oSelf *AppUserModel) ShowOneByName(sName string) (*domain.AppUser, error) {
	var oAppUser domain.AppUser

	if oErr := oSelf.
		DB.
		WithContext(oSelf.Context).
		Where("name = ?", sName).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oAppUser).
		Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, errors.New("資料不存在")
		}
		return nil, oErr
	}

	return &oAppUser, nil
}

func (oSelf *AppUserModel) ShowOneById(iId uint) (*domain.AppUser, error) {
	var oAppUser domain.AppUser

	if oErr := oSelf.
		DB.
		WithContext(oSelf.Context).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oAppUser, iId).
		Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, errors.New("資料不存在")
		}
		return nil, oErr
	}

	return &oAppUser, nil
}

func (oSelf *AppUserModel) AddOne(oValue *domain.AppUserValue) error {
	oColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).Model(&domain.AppUser{}).Create(oColumns)
	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("新增0筆")
	}

	return nil
}
