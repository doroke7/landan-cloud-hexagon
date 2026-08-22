package postgresql

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	domain "example/internal/domain"
	postgresqlBase "example/internal/output/application/postgresql"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkg "example/pkg"
)

type TableModel struct {
	*postgresqlBase.AbstractPostgresql
}

func NewTableModel(oAbstractModel *postgresqlBase.AbstractPostgresql) outputPortAnyModel.TableModel {
	return &TableModel{
		AbstractPostgresql: oAbstractModel,
	}
}

func (oSelf *TableModel) ShowOneById(iId uint) (*domain.Table, error) {
	var oTableRow domain.TableRow

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Preload("Game").
		Preload("Game.GameType").
		Model(&oTableRow).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oTableRow, iId).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return domain.TableRowToTable(&oTableRow), nil
}

func (oSelf *TableModel) EditOneById(oTable *domain.TableValue, iId uint) (bool, error) {
	var oTableRow domain.TableRow

	oColumns, oErr := pkg.StructToMap(oTable)
	if oErr != nil {
		return false, oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oTableRow).
		Where("id = ?", iId).
		UpdateColumns(oColumns)

	if oResult.Error != nil {
		return false, oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return false, errors.New("更新0筆")
	}

	return true, nil
}

func (oSelf *TableModel) RemoveOneById(iId uint) (bool, error) {
	var oTableRow domain.TableRow

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oTableRow).
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

func (oSelf *TableModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Table, error) {
	aWheres := oSelf.AbstractPostgresql.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractPostgresql.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractPostgresql.PaginationToLimit(oPagination)

	var aTableRows []*domain.TableRow
	var oTableRow domain.TableRow

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Preload("Game").
		Preload("Game.GameType").
		Model(&oTableRow).
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
		Find(&aTableRows).Error; oErr != nil {
		return nil, oErr
	}

	aTables := make([]*domain.Table, len(aTableRows))
	for i, oTableRow := range aTableRows {
		aTables[i] = domain.TableRowToTable(oTableRow)
	}

	return aTables, nil
}

func (oSelf *TableModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	aWheres := oSelf.AbstractPostgresql.FiltersToWheres(aFilters)

	var iTotal int64
	var oTableRow domain.TableRow

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&oTableRow)

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	if oErr := oQuery.
		Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}

func (oSelf *TableModel) AddOne(oTable *domain.TableValue) (bool, error) {
	var oTableRow domain.TableRow

	oColumns, oErr := pkg.StructToMap(oTable)
	if oErr != nil {
		return false, oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oTableRow).
		Create(oColumns)

	if oResult.Error != nil {
		return false, oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return false, errors.New("新增0筆")
	}

	return true, nil
}
