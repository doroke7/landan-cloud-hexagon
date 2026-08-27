package outputApplicationPostgresqlModel

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	domain "example/internal/domain"
	postgresqlBase "example/internal/output/application/postgresql"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type GameTypeModel struct {
	*postgresqlBase.AbstractPostgresql
}

func NewGameTypeModel(oAbstractModel *postgresqlBase.AbstractPostgresql) outputPortAnyModel.GameTypeModel {
	return &GameTypeModel{
		AbstractPostgresql: oAbstractModel,
	}
}

func (oSelf *GameTypeModel) ShowOneById(iId uint) (*domain.GameType, error) {
	var oGameTypeRow domain.GameTypeRow

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Model(&oGameTypeRow).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oGameTypeRow, iId).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return domain.GameTypeRowToGameType(&oGameTypeRow), nil
}

func (oSelf *GameTypeModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, error) {
	aWheres := oSelf.AbstractPostgresql.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractPostgresql.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractPostgresql.PaginationToLimit(oPagination)

	var aGameTypeRows []*domain.GameTypeRow
	var oGameTypeRow domain.GameTypeRow

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&oGameTypeRow).
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
		Find(&aGameTypeRows).Error; oErr != nil {
		return nil, oErr
	}

	aGameTypes := make([]*domain.GameType, len(aGameTypeRows))
	for i, oGameTypeRow := range aGameTypeRows {
		aGameTypes[i] = domain.GameTypeRowToGameType(oGameTypeRow)
	}

	return aGameTypes, nil
}

func (oSelf *GameTypeModel) AddOne(oValue *domain.GameTypeValue) (bool, error) {
	var oGameTypeRow domain.GameTypeRow

	oGameType, _ := pkgUtility.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oGameTypeRow).
		Create(oGameType)

	if oResult.Error != nil {
		return false, oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return false, errors.New("新增0筆")
	}

	return true, nil
}

func (oSelf *GameTypeModel) EditOneById(oValue *domain.GameTypeValue, iId uint) (bool, error) {
	var oGameTypeRow domain.GameTypeRow

	oGameType, _ := pkgUtility.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oGameTypeRow).
		Where("id = ?", iId).
		UpdateColumns(oGameType)

	if oResult.Error != nil {
		return false, oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return false, errors.New("更新0筆")
	}

	return true, nil
}

func (oSelf *GameTypeModel) RemoveOneById(iId uint) (bool, error) {
	var oGameTypeRow domain.GameTypeRow

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oGameTypeRow).
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

func (oSelf *GameTypeModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aWheres := oSelf.AbstractPostgresql.FiltersToWheres(aFilters)

	var iTotal int64
	var oGameTypeRow domain.GameTypeRow

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&oGameTypeRow)

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	if oErr := oQuery.
		Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}
