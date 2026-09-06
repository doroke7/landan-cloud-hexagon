package outputApplicationMysqlModel

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type TableRecordLogModel struct {
	*AbstractModel
}

func NewTableRecordLogModel(oAbstractModel *AbstractModel) outputPortAnyModel.TableRecordLogModel {
	return &TableRecordLogModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *TableRecordLogModel) ShowOneById(iId uint64) (*domain.TableRecordLog, error) {
	var oTableRecordLog domain.TableRecordLog

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oTableRecordLog).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oTableRecordLog, iId)

	if oResult.Error != nil {
		if errors.Is(oResult.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oResult.Error
	}

	return &oTableRecordLog, nil
}

func (oSelf *TableRecordLogModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.TableRecordLog, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractMysql.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractMysql.PaginationToLimit(oPagination)

	var aTableRecordLogs []*domain.TableRecordLog

	oQuery := oSelf.DB.WithContext(oSelf.Context).
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

	oResult := oQuery.
		Limit(int(*oLimit.Count)).
		Offset(int(*oLimit.Offset)).
		Find(&aTableRecordLogs)

	if oResult.Error != nil {
		return nil, oResult.Error
	}

	return aTableRecordLogs, nil
}

func (oSelf *TableRecordLogModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)

	var iTotal int64
	oQuery := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.TableRecordLog{})

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
		return errors.New("0 rows inserted")
	}

	return nil
}
