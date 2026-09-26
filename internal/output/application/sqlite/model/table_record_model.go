package outputApplicationSqliteModel

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	domain "example/internal/domain"
	outputApplicationSqlite "example/internal/output/application/sqlite"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type TableRecordModel struct {
	*outputApplicationSqlite.AbstractSqlite
}

func NewTableRecordModel(oAbstractSqlite *outputApplicationSqlite.AbstractSqlite) outputPortAnyModel.TableRecordModel {
	return &TableRecordModel{
		AbstractSqlite: oAbstractSqlite,
	}
}

func (oSelf *TableRecordModel) ShowOneById(iId uint64) (*domain.TableRecord, error) {
	var oTableRecord domain.TableRecord

	oErr := oSelf.DB.WithContext(oSelf.Context).
		Model(&oTableRecord).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oTableRecord, iId).Error

	if oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return &oTableRecord, nil
}

func (oSelf *TableRecordModel) EditOneById(oTableRecord *domain.TableRecordVariable, iId uint64) error {
	oColumns, oErr := pkgUtility.StructToMap(oTableRecord)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.TableRecord{}).
		Where("id = ?", iId).
		UpdateColumns(oColumns)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		oZeroRowsError := errors.New("0 rows updated")

		return oZeroRowsError
	}

	return nil
}

func (oSelf *TableRecordModel) RemoveOneById(iId uint64) error {
	oNow := time.Now()

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.TableRecord{}).
		Where("id = ?", iId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		UpdateColumn("deleted_at", oNow)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		oZeroRowsError := errors.New("0 rows deleted")

		return oZeroRowsError
	}

	return nil
}

func (oSelf *TableRecordModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.TableRecord, error) {
	aWheres := oSelf.AbstractSqlite.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractSqlite.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractSqlite.PaginationToLimit(oPagination)

	var aTableRecords []*domain.TableRecord

	oQuery := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.TableRecord{}).
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

	oErr := oQuery.
		Limit(int(*oLimit.Count)).
		Offset(int(*oLimit.Offset)).
		Find(&aTableRecords).Error

	if oErr != nil {
		return nil, oErr
	}

	return aTableRecords, nil
}

func (oSelf *TableRecordModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aWheres := oSelf.AbstractSqlite.FiltersToWheres(aFilters)

	var iTotal int64
	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&domain.TableRecord{})

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	oErr := oQuery.Count(&iTotal).Error

	if oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}

func (oSelf *TableRecordModel) AddOne(oTableRecord *domain.TableRecordVariable) error {
	oColumns, oErr := pkgUtility.StructToMap(oTableRecord)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.TableRecord{}).
		Create(oColumns)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		oZeroRowsError := errors.New("0 rows inserted")

		return oZeroRowsError
	}

	return nil
}
