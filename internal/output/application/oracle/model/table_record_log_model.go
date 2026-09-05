package outputApplicationOracleModel

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	domain "example/internal/domain"
	oracleBase "example/internal/output/application/oracle"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type TableRecordLogModel struct {
	*oracleBase.AbstractOracle
}

func NewTableRecordLogModel(oAbstractModel *oracleBase.AbstractOracle) outputPortAnyModel.TableRecordLogModel {
	return &TableRecordLogModel{
		AbstractOracle: oAbstractModel,
	}
}

func (oSelf *TableRecordLogModel) ShowOneById(iId uint) (*domain.TableRecordLog, error) {
	var oTableRecordLog domain.TableRecordLog

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Model(&oTableRecordLog).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oTableRecordLog, iId).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return &oTableRecordLog, nil
}

func (oSelf *TableRecordLogModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.TableRecordLog, error) {
	aWheres := oSelf.AbstractOracle.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractOracle.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractOracle.PaginationToLimit(oPagination)

	var aTableRecordLogs []*domain.TableRecordLog

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&domain.TableRecordLog{}).
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
		Find(&aTableRecordLogs).Error; oErr != nil {
		return nil, oErr
	}

	return aTableRecordLogs, nil
}

func (oSelf *TableRecordLogModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint, error) {
	aWheres := oSelf.AbstractOracle.FiltersToWheres(aFilters)

	var iTotal int64
	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&domain.TableRecordLog{})

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	if oErr := oQuery.
		Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint(iTotal), nil
}

func (oSelf *TableRecordLogModel) AddOne(oTableRecordLog *domain.TableRecordLogValue) error {
	oColumns, oErr := pkgUtility.StructToMap(oTableRecordLog)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.TableRecordLog{}).
		Create(oColumns)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("新增0筆")
	}

	return nil
}
