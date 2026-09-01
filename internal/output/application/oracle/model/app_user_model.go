package outputApplicationOracleModel

import (
	"errors"

	"gorm.io/gorm"

	domain "example/internal/domain"
	oracleBase "example/internal/output/application/oracle"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgUtility "example/pkg/utility"
)

type AppUserModel struct {
	*oracleBase.AbstractOracle
}

func NewAppUserModel(oAbstractModel *oracleBase.AbstractOracle) outputPortAnyModel.AppUserModel {
	return &AppUserModel{
		AbstractOracle: oAbstractModel,
	}
}

func (oSelf *AppUserModel) IncreaseBalance(id uint, amount uint) (bool, error) {
	if err := oSelf.DB.WithContext(oSelf.Context).Model(&domain.AppUser{}).
		Where("id = ?", id).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (oSelf *AppUserModel) ShowOneByName(sName string) (*domain.AppUser, error) {
	var oAppUser domain.AppUser

	if oErr := oSelf.DB.WithContext(oSelf.Context).Where("name = ?", sName).First(&oAppUser).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, errors.New("資料不存在")
		}
		return nil, oErr
	}

	return &oAppUser, nil
}

func (oSelf *AppUserModel) ShowOneById(iId uint) (*domain.AppUser, error) {
	var oAppUser domain.AppUser

	if oErr := oSelf.DB.WithContext(oSelf.Context).First(&oAppUser, iId).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, errors.New("資料不存在")
		}
		return nil, oErr
	}

	return &oAppUser, nil
}

func (oSelf *AppUserModel) AddOne(oValue *domain.AppUserValue) (bool, error) {
	oColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return false, oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).Model(&domain.AppUser{}).Create(oColumns)
	if oResult.Error != nil {
		return false, oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return false, errors.New("新增0筆")
	}

	return true, nil
}
