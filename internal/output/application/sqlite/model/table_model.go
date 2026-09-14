package outputApplicationSqliteModel

import (
	"errors"
	"time"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type TableModel struct {
	*AbstractModel
}

func NewTableModel(oAbstractModel *AbstractModel) outputPortAnyModel.TableModel {
	return &TableModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *TableModel) EditOneById(oTable *domain.TableVariable, iId uint64) error {
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

	if oResult.RowsAffected == 0 {
		oZeroRowsError := errors.New("0 rows updated")

		return oZeroRowsError
	}

	return nil
}

func (oSelf *TableModel) RemoveOneById(iId uint64) error {
	oNow := time.Now()

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Table{}).
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

func (oSelf *TableModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aWheres := oSelf.AbstractSqlite.FiltersToWheres(aFilters)

	var iTotal int64
	oQuery := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.Table{})

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	oErr := oQuery.Count(&iTotal).Error

	if oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}

func (oSelf *TableModel) AddOne(oTable *domain.TableVariable) error {

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
		oZeroRowsError := errors.New("0 rows inserted")

		return oZeroRowsError
	}

	return nil
}
