package outputApplicationMysqlModel

import (
	"errors"
	"time"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type GameTypeModel struct {
	*AbstractModel
}

func NewGameTypeModel(oAbstractModel *AbstractModel) outputPortAnyModel.GameTypeModel {
	return &GameTypeModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *GameTypeModel) TotalByParentId(iParentId uint64) (uint64, error) {
	var iTotal int64

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.GameType{}).
		Where("parent_id = ?", iParentId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Count(&iTotal)

	if oResult.Error != nil {
		return 0, oResult.Error
	}

	return uint64(iTotal), nil
}

func (oSelf *GameTypeModel) AddOne(oValue *domain.GameTypeValue) error {
	oGameType, _ := pkgUtility.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.GameType{}).
		Create(oGameType)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("0 rows inserted")
	}

	return nil
}

func (oSelf *GameTypeModel) EditOneById(oValue *domain.GameTypeValue, iId uint64) error {
	oGameType, _ := pkgUtility.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.GameType{}).
		Where("id = ?", iId).
		UpdateColumns(oGameType)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("0 rows updated")
	}

	return nil
}

func (oSelf *GameTypeModel) RemoveOneById(iId uint64) error {

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.GameType{}).
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

func (oSelf *GameTypeModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)

	var iTotal int64
	oQuery := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.GameType{})

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	oResult := oQuery.
		Count(&iTotal)

	if oResult.Error != nil {
		return 0, oResult.Error
	}

	return uint64(iTotal), nil
}
