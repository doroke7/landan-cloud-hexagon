package outputApplicationMysqlModel

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	domain "example/internal/domain"
	outputApplicationMysql "example/internal/output/application/mysql"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type GameTypeModel struct {
	*outputApplicationMysql.AbstractMysql
}

func NewGameTypeModel(oAbstractModel *outputApplicationMysql.AbstractMysql) outputPortAnyModel.GameTypeModel {
	return &GameTypeModel{
		AbstractMysql: oAbstractModel,
	}
}

func (oSelf *GameTypeModel) ShowOneById(iId uint64) (*domain.GameType, error) {
	var oGameType domain.GameType

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Model(&oGameType).
		Preload("Parent").
		Preload("Children").
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oGameType, iId).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return &oGameType, nil
}

func (oSelf *GameTypeModel) ShowOnes() ([]*domain.GameType, error) {
	var aGameTypes []*domain.GameType

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.GameType{}).
		Preload("Parent").
		Preload("Children").
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Find(&aGameTypes).Error; oErr != nil {
		return nil, oErr
	}

	return aGameTypes, nil
}

// ShowOnesByParentId 撈出指定父類型底下、尚未刪除的子類型（給刪除前的擋關檢查用）。
func (oSelf *GameTypeModel) TotalByParentId(iParentId uint64) (uint64, error) {
	var iTotal int64

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.GameType{}).
		Where("parent_id = ?", iParentId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}

func (oSelf *GameTypeModel) ShowOnesByParentId(iParentId uint64) ([]*domain.GameType, error) {
	var aGameTypes []*domain.GameType

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.GameType{}).
		Preload("Parent").
		Preload("Children").
		Where("parent_id = ?", iParentId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Find(&aGameTypes).Error; oErr != nil {
		return nil, oErr
	}

	return aGameTypes, nil
}

func (oSelf *GameTypeModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractMysql.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractMysql.PaginationToLimit(oPagination)

	var aGameTypes []*domain.GameType

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Preload("Parent").
		Preload("Children").
		Model(&domain.GameType{}).
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
		Find(&aGameTypes).Error; oErr != nil {
		return nil, oErr
	}

	return aGameTypes, nil
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
	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&domain.GameType{})

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	if oErr := oQuery.
		Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}
