package gaussdb

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	domain "example/internal/domain"
	gaussdbBase "example/internal/output/application/gaussdb"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkg "example/pkg"
)

type TableRecordModel struct {
	*gaussdbBase.AbstractGaussdb
}

func NewTableRecordModel(oAbstractModel *gaussdbBase.AbstractGaussdb) outputPortAnyModel.TableRecordModel {
	return &TableRecordModel{
		AbstractGaussdb: oAbstractModel,
	}
}

func (oSelf *TableRecordModel) ShowOneById(iId uint) (*domain.TableRecord, error) {
	var oTableRecordRow domain.TableRecordRow

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Model(&oTableRecordRow).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oTableRecordRow, iId).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return domain.TableRecordRowToTableRecord(&oTableRecordRow), nil
}

func (oSelf *TableRecordModel) EditOneById(oTableRecord *domain.TableRecordValue, iId uint) (bool, error) {
	var oTableRecordRow domain.TableRecordRow

	oColumns, oErr := pkg.StructToMap(oTableRecord)
	if oErr != nil {
		return false, oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oTableRecordRow).
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

func (oSelf *TableRecordModel) RemoveOneById(iId uint) (bool, error) {
	var oTableRecordRow domain.TableRecordRow

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oTableRecordRow).
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

func (oSelf *TableRecordModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.TableRecord, error) {
	aWheres := oSelf.AbstractGaussdb.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractGaussdb.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractGaussdb.PaginationToLimit(oPagination)

	var aTableRecordRows []*domain.TableRecordRow
	var oTableRecordRow domain.TableRecordRow

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&oTableRecordRow).
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
		Find(&aTableRecordRows).Error; oErr != nil {
		return nil, oErr
	}

	aTableRecords := make([]*domain.TableRecord, len(aTableRecordRows))
	for i, oTableRecordRow := range aTableRecordRows {
		aTableRecords[i] = domain.TableRecordRowToTableRecord(oTableRecordRow)
	}

	return aTableRecords, nil
}

func (oSelf *TableRecordModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	aWheres := oSelf.AbstractGaussdb.FiltersToWheres(aFilters)

	var iTotal int64
	var oTableRecordRow domain.TableRecordRow

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&oTableRecordRow)

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	if oErr := oQuery.
		Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}

func (oSelf *TableRecordModel) AddOne(oTableRecord *domain.TableRecordValue) (bool, error) {
	var oTableRecordRow domain.TableRecordRow

	oColumns, oErr := pkg.StructToMap(oTableRecord)
	if oErr != nil {
		return false, oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oTableRecordRow).
		Create(oColumns)

	if oResult.Error != nil {
		return false, oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return false, errors.New("新增0筆")
	}

	return true, nil
}
