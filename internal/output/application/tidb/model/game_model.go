package outputApplicationTidbModel

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	domain "example/internal/domain"
	tidbBase "example/internal/output/application/tidb"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type GameModel struct {
	*tidbBase.AbstractTidb
}

func NewGameModel(oAbstractModel *tidbBase.AbstractTidb) outputPortAnyModel.GameModel {
	return &GameModel{
		AbstractTidb: oAbstractModel,
	}
}

func (oSelf *GameModel) ShowOneById(iId uint) (*domain.Game, error) {
	var oGame domain.Game

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Preload("GameType").
		Preload("GameType.Parent").
		Model(&oGame).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oGame, iId).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return &oGame, nil
}

func (oSelf *GameModel) ShowOneByKey(sKey string) (*domain.Game, error) {
	var oGame domain.Game

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Preload("GameType").
		Preload("GameType.Parent").
		Model(&oGame).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Where(map[string]any{"key": sKey}).
		First(&oGame).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return &oGame, nil
}

func (oSelf *GameModel) ShowOnesByGameTypeId(iGameTypeId uint) ([]*domain.Game, error) {
	var aGames []*domain.Game

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Preload("GameType").
		Preload("GameType.Parent").
		Model(&domain.Game{}).
		Where("game_type_id = ?", iGameTypeId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Find(&aGames).Error; oErr != nil {
		return nil, oErr
	}

	return aGames, nil
}

func (oSelf *GameModel) TotalByGameTypeId(iGameTypeId uint) (uint, error) {
	var iTotal int64

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Game{}).
		Where("game_type_id = ?", iGameTypeId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint(iTotal), nil
}

func (oSelf *GameModel) ShowOnesByFiltersWithOrdersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, error) {
	aWheres := oSelf.AbstractTidb.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractTidb.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractTidb.PaginationToLimit(oPagination)

	var aGames []*domain.Game

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Preload("GameType").
		Preload("GameType.Parent").
		Model(&domain.Game{}).
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
		Find(&aGames).Error; oErr != nil {
		return nil, oErr
	}

	return aGames, nil
}

func (oSelf *GameModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint, error) {
	aWheres := oSelf.AbstractTidb.FiltersToWheres(aFilters)

	var iTotal int64
	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&domain.Game{})

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	// Go 很常先在外面宣告變數，然後把它的 pointer 傳進函數，讓函數可以直接修改原本的變數。
	if oErr := oQuery.
		Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint(iTotal), nil

}

func (oSelf *GameModel) AddOne(oValue *domain.GameValue) error {
	oGame, _ := pkgUtility.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Game{}).
		Create(oGame)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("新增0筆")

	}

	return nil
}

func (oSelf *GameModel) EditOneById(oValue *domain.GameValue, iId uint) error {
	oGame, _ := pkgUtility.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Game{}).
		Where("id = ?", iId).
		UpdateColumns(oGame)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("更新0筆")

	}

	return nil
}

func (oSelf *GameModel) RemoveOneById(iId uint) error {
	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Game{}).
		Where("id = ?", iId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		UpdateColumn("deleted_at", time.Now())

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("刪除0筆")

	}

	return nil
}
