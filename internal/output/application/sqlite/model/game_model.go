package outputApplicationSqliteModel

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	domain "example/internal/domain"
	sqliteBase "example/internal/output/application/sqlite"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type GameModel struct {
	*sqliteBase.AbstractSqlite
}

func NewGameModel(oAbstractModel *sqliteBase.AbstractSqlite) outputPortAnyModel.GameModel {
	return &GameModel{
		AbstractSqlite: oAbstractModel,
	}
}

func (oSelf *GameModel) ShowOneById(iId uint) (*domain.Game, error) {
	var oGameRow domain.GameRow

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Preload("GameType").
		Model(&oGameRow).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oGameRow, iId).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return domain.GameRowToGame(&oGameRow), nil
}

func (oSelf *GameModel) ShowOnesByFiltersWithOrdersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, error) {
	aWheres := oSelf.AbstractSqlite.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractSqlite.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractSqlite.PaginationToLimit(oPagination)

	var aGameRows []*domain.GameRow
	var oGameRow domain.GameRow

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Preload("GameType").
		Model(&oGameRow).
		Where("deleted_at = ?", "2038-01-19 03:14:07")

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	for _, oOrder := range aOrders {
		if oOrder == nil || oOrder.Field == nil {
			continue
		}

		sDirection := "ASC"
		if oOrder.Value != nil && strings.EqualFold(*oOrder.Value, "desc") {
			sDirection = "DESC"
		}

		oQuery = oQuery.Order(*oOrder.Field + " " + sDirection)
	}

	if oErr := oQuery.
		Limit(int(*oLimit.Count)).
		Offset(int(*oLimit.Offset)).
		Find(&aGameRows).Error; oErr != nil {
		return nil, oErr
	}

	aGames := make([]*domain.Game, len(aGameRows))
	for i, oGameRow := range aGameRows {
		aGames[i] = domain.GameRowToGame(oGameRow)
	}

	return aGames, nil
}

func (oSelf *GameModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aWheres := oSelf.AbstractSqlite.FiltersToWheres(aFilters)

	var iTotal int64
	var oGameRow domain.GameRow

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&oGameRow)

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	if oErr := oQuery.
		Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}

func (oSelf *GameModel) AddOne(oValue *domain.GameValue) (bool, error) {
	var oGameRow domain.GameRow

	oGame, _ := pkgUtility.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oGameRow).
		Create(oGame)

	if oResult.Error != nil {
		return false, oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return false, errors.New("新增0筆")
	}

	return true, nil
}

func (oSelf *GameModel) EditOneById(oValue *domain.GameValue, iId uint) (bool, error) {
	var oGameRow domain.GameRow

	oGame, _ := pkgUtility.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oGameRow).
		Where("id = ?", iId).
		UpdateColumns(oGame)

	if oResult.Error != nil {
		return false, oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return false, errors.New("更新0筆")
	}

	return true, nil
}

func (oSelf *GameModel) RemoveOneById(iId uint) (bool, error) {
	var oGameRow domain.GameRow

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oGameRow).
		Where("id = ?", iId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		UpdateColumn("deleted_at", time.Now())

	if oResult.Error != nil {
		return false, oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return false, errors.New("刪除0筆")
	}

	return true, nil
}
