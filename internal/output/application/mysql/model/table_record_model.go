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

type TableRecordModel struct {
	*outputApplicationMysql.AbstractMysql
}

func NewTableRecordModel(oAbstractMysql *outputApplicationMysql.AbstractMysql) outputPortAnyModel.TableRecordModel {
	return &TableRecordModel{
		AbstractMysql: oAbstractMysql,
	}
}

func (oSelf *TableRecordModel) ShowOneById(iId uint64) (*domain.TableRecord, error) {
	var oTableRecord domain.TableRecord

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oTableRecord).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oTableRecord, iId)

	if oResult.Error != nil {
		if errors.Is(oResult.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oResult.Error
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
		return errors.New("0 rows updated")
	}

	return nil
}

func (oSelf *TableRecordModel) RemoveOneById(iId uint64) error {
	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.TableRecord{}).
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

func (oSelf *TableRecordModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.TableRecord, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractMysql.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractMysql.PaginationToLimit(oPagination)

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

	oResult := oQuery.
		Limit(int(*oLimit.Count)).
		Offset(int(*oLimit.Offset)).
		Find(&aTableRecords)

	if oResult.Error != nil {
		return nil, oResult.Error
	}

	return aTableRecords, nil
}

func (oSelf *TableRecordModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)

	var iTotal int64
	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&domain.TableRecord{})

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
		return errors.New("0 rows inserted")
	}

	return nil
}
