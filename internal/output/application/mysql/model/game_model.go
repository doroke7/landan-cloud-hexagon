package outputApplicationMysqlModel

import (
	"errors"
	"time"

	"gorm.io/gorm"

	domain "example/internal/domain"
	outputApplicationMysql "example/internal/output/application/mysql"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type GameModel struct {
	*outputApplicationMysql.AbstractMysql
}

func NewGameModel(oAbstractMysql *outputApplicationMysql.AbstractMysql) outputPortAnyModel.GameModel {
	return &GameModel{
		AbstractMysql: oAbstractMysql,
	}
}

func (oSelf *GameModel) ShowOneByKey(sKey string) (*domain.Game, error) {
	var oGame domain.Game

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oGame).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Where(map[string]any{"key": sKey}).
		First(&oGame)

	if oResult.Error != nil {
		if errors.Is(oResult.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oResult.Error
	}

	return &oGame, nil
}

func (oSelf *GameModel) TotalByGameTypeId(iGameTypeId uint64) (uint, error) {
	var iTotal int64

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Game{}).
		Where("game_type_id = ?", iGameTypeId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Count(&iTotal)

	if oResult.Error != nil {
		return 0, oResult.Error
	}

	return uint(iTotal), nil
}

func (oSelf *GameModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)

	var iTotal int64
	oQuery := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Game{})

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	// Go 很常先在外面宣告變數，然後把它的 pointer 傳進函數，讓函數可以直接修改原本的變數。
	oResult := oQuery.
		Count(&iTotal)

	if oResult.Error != nil {
		return 0, oResult.Error
	}

	return uint64(iTotal), nil

}

func (oSelf *GameModel) AddOne(oValue *domain.GameVariable) error {
	oGame, _ := pkgUtility.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Game{}).
		Create(oGame)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("0 rows inserted")

	}

	return nil
}

func (oSelf *GameModel) EditOneById(oValue *domain.GameVariable, iId uint64) error {
	oGame, _ := pkgUtility.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Game{}).
		Where("id = ?", iId).
		UpdateColumns(oGame)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("0 rows updated")

	}

	return nil
}

func (oSelf *GameModel) RemoveOneById(iId uint64) error {
	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Game{}).
		Where("id = ?", iId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		UpdateColumn("deleted_at", time.Now())

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("0 rows deleted")

	}

	return nil
}
