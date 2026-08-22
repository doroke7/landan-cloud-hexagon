package mysql

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkg "example/pkg"
)

type TableRecordLogModel struct {
	*AbstractModel
}

func NewTableRecordLogModel(oAbstractModel *AbstractModel) outputPortAnyModel.TableRecordLogModel {
	return &TableRecordLogModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *TableRecordLogModel) ShowOneById(iId uint) (*domain.TableRecordLog, error) {
	var oTableRecordLogRow domain.TableRecordLogRow

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Model(&oTableRecordLogRow).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oTableRecordLogRow, iId).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return domain.TableRecordLogRowToTableRecordLog(&oTableRecordLogRow), nil
}

func (oSelf *TableRecordLogModel) ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.TableRecordLog, error) {
	var aTableRecordLogRows []*domain.TableRecordLogRow
	var oTableRecordLogRow domain.TableRecordLogRow

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&oTableRecordLogRow).
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
		Find(&aTableRecordLogRows).Error; oErr != nil {
		return nil, oErr
	}

	aTableRecordLogs := make([]*domain.TableRecordLog, len(aTableRecordLogRows))
	for i, oTableRecordLogRow := range aTableRecordLogRows {
		aTableRecordLogs[i] = domain.TableRecordLogRowToTableRecordLog(oTableRecordLogRow)
	}

	return aTableRecordLogs, nil
}

func (oSelf *TableRecordLogModel) TotalByWheres(aWheres []*pkg.Where) (uint64, error) {
	var iTotal int64
	var oTableRecordLogRow domain.TableRecordLogRow

	oQuery := oSelf.
		DB.
		WithContext(oSelf.Context).
		Model(&oTableRecordLogRow)

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	if oErr := oQuery.
		Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}

func (oSelf *TableRecordLogModel) AddOne(oTableRecordLog *domain.TableRecordLogValue) (bool, error) {
	var oTableRecordLogRow domain.TableRecordLogRow

	oColumns, oErr := pkg.StructToMap(oTableRecordLog)
	if oErr != nil {
		return false, oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oTableRecordLogRow).
		Create(oColumns)

	if oResult.Error != nil {
		return false, oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return false, errors.New("新增0筆")
	}

	return true, nil
}
