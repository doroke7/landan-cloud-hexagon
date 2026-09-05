package outputApplicationClickhouseModel

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	domain "example/internal/domain"
	clickhouseBase "example/internal/output/application/clickhouse"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type TableModel struct {
	*clickhouseBase.AbstractClickhouse
}

func NewTableModel(oAbstractModel *clickhouseBase.AbstractClickhouse) outputPortAnyModel.TableModel {
	return &TableModel{
		AbstractClickhouse: oAbstractModel,
	}
}

func (oSelf *TableModel) ShowOneById(iId uint) (*domain.Table, error) {
	var oTable domain.Table

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Preload("Game").
		Preload("Game.GameType").
		Model(&oTable).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oTable, iId).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return &oTable, nil
}

// EditOneById 底層送出的是 ClickHouse 的 ALTER TABLE ... UPDATE mutation，
// 非同步執行，RowsAffected 通常會回 0，不能拿來判斷是否真的更新成功。
func (oSelf *TableModel) EditOneById(oTable *domain.TableValue, iId uint64) error {
	oColumns, oErr := pkgUtility.StructToMap(oTable)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Table{}).
		Where("id = ?", iId).
		UpdateColumns(oColumns)

	if oResult.Error != nil {
		return oResult.Error
	}

	return nil
}

// RemoveOneById 底層送出的是 ClickHouse 的 ALTER TABLE ... UPDATE mutation（軟刪除），
// 非同步執行，RowsAffected 通常會回 0，不能拿來判斷是否真的刪除成功。
func (oSelf *TableModel) RemoveOneById(iId uint64) error {
	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Table{}).
		Where("id = ?", iId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		UpdateColumn("deleted_at", time.Now())

	if oResult.Error != nil {
		return oResult.Error
	}

	return nil
}

func (oSelf *TableModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, error) {
	aWheres := oSelf.AbstractClickhouse.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractClickhouse.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractClickhouse.PaginationToLimit(oPagination)

	var aTables []*domain.Table

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Preload("Game").
		Preload("Game.GameType").
		Model(&domain.Table{}).
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
		Find(&aTables).Error; oErr != nil {
		return nil, oErr
	}

	return aTables, nil
}

func (oSelf *TableModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aWheres := oSelf.AbstractClickhouse.FiltersToWheres(aFilters)

	var iTotal int64
	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&domain.Table{})

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	if oErr := oQuery.
		Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}

func (oSelf *TableModel) AddOne(oTable *domain.TableValue) error {
	oColumns, oErr := pkgUtility.StructToMap(oTable)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Table{}).
		Create(oColumns)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("新增0筆")
	}

	return nil
}
