package clickhouse

import (
	"context"
	"strings"

	pkg "example/pkg"

	"gorm.io/gorm"
)

var oOperatorMap = map[string]string{
	"eq":          "=",
	"ne":          "!=",
	"gt":          ">",
	"gte":         ">=",
	"lt":          "<",
	"lte":         "<=",
	"contains":    "LIKE",
	"notContains": "NOT LIKE",
	"startsWith":  "LIKE",
	"endsWith":    "LIKE",
	"in":          "IN",
	"notIn":       "NOT IN",
	"between":     "BETWEEN",
}

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractRepository 做法一致。
type AbstractClickhouse struct {
	DB      *gorm.DB
	Context context.Context
	*pkg.Aop
}

func NewAbstractClickhouse(oContext context.Context, oDb *gorm.DB, oAop *pkg.Aop) *AbstractClickhouse {

	return &AbstractClickhouse{
		DB:      oDb,
		Context: oContext,
		Aop:     oAop,
	}
}

func (oSelf *AbstractClickhouse) FiltersToWheres(aFilters []*pkg.Filter) []*pkg.ClickhouseWhere {
	aWheres := make([]*pkg.ClickhouseWhere, 0, len(aFilters))

	for _, oFilter := range aFilters {
		if oFilter == nil || oFilter.Field == nil {
			continue
		}

		sOperator := "="
		oValue := oFilter.Value

		if oFilter.Operator != nil {
			if sSql, bOk := oOperatorMap[*oFilter.Operator]; bOk {
				sOperator = sSql
			}

			if sValue, bOk := oValue.(string); bOk {
				switch *oFilter.Operator {
				case "contains", "notContains":
					oValue = "%" + sValue + "%"
				case "startsWith":
					oValue = sValue + "%"
				case "endsWith":
					oValue = "%" + sValue
				}
			}
		}

		// ClickHouse 的識別字用反引號，跟 mysql 一樣。
		sField := "`" + *oFilter.Field + "`"
		aWheres = append(aWheres, &pkg.ClickhouseWhere{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oValue,
		})
	}

	return aWheres
}

func (oSelf *AbstractClickhouse) SortersToOrders(aSorters []*pkg.Sorter) []*pkg.ClickhouseOrder {
	aOrders := make([]*pkg.ClickhouseOrder, 0, len(aSorters))

	for _, oSorter := range aSorters {
		if oSorter == nil || oSorter.Field == nil {
			continue
		}

		sDirection := "asc"
		if oSorter.Order != nil && strings.EqualFold(*oSorter.Order, "desc") {
			sDirection = "desc"
		}

		sField := "`" + *oSorter.Field + "`"
		aOrders = append(aOrders, &pkg.ClickhouseOrder{
			Field: &sField,
			Value: &sDirection,
		})
	}

	return aOrders
}

func (oSelf *AbstractClickhouse) PaginationToLimit(oPagination *pkg.Pagination) *pkg.ClickhouseLimit {
	iSize := uint(10)
	iPage := uint(1)

	if oPagination != nil && oPagination.Size != nil && *oPagination.Size != 0 {
		iSize = *oPagination.Size
	}

	if oPagination != nil && oPagination.Page != nil && *oPagination.Page != 0 {
		iPage = *oPagination.Page
	}

	iOffset := (iPage - 1) * iSize

	return &pkg.ClickhouseLimit{
		Offset: &iOffset,
		Count:  &iSize,
	}
}
