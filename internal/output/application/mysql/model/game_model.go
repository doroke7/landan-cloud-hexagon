package mysql

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkg "example/pkg"
)

type GameModel struct {
	*AbstractModel
}

func NewGameModel(oAbstractModel *AbstractModel) outputPortAnyModel.GameModel {
	return &GameModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *GameModel) ShowOneById(iId uint) (*domain.Game, error) {
	var oGame domain.Game

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Preload("GameType").
		Model(&domain.Game{}).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oGame, iId).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return &oGame, nil
}

func (oSelf *GameModel) ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Game, error) {
	var aGames []*domain.Game

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Preload("GameType").
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

func (oSelf *GameModel) AddOne(oValue *domain.GameValue) (bool, error) {

	oGame, _ := pkg.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Game{}).
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
	oGame, _ := pkg.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Game{}).
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

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Game{}).
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

func (oSelf *GameModel) TotalByWheres(aWheres []*pkg.Where) (uint64, error) {
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

	return uint64(iTotal), nil
}
