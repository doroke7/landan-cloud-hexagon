package outputApplicationClickhouse

import (
	"context"
	"strings"

	pkgCache "example/pkg/cache"
	pkgClickhouse "example/pkg/clickhouse"
	pkgInput "example/pkg/input"

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

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkgCache.Aop、cache/memory 的
// AbstractRepository 做法一致。
type AbstractClickhouse struct {
	DB      *gorm.DB
	Context context.Context
	*pkgCache.Aop
}

func NewAbstractClickhouse(oContext context.Context, oDb *gorm.DB, oAop *pkgCache.Aop) *AbstractClickhouse {

	return &AbstractClickhouse{
		DB:      oDb,
		Context: oContext,
		Aop:     oAop,
	}
}

func (oSelf *AbstractClickhouse) FiltersToWheres(aFilters []*pkgInput.Filter) []*pkgClickhouse.ClickhouseWhere {
	aWheres := make([]*pkgClickhouse.ClickhouseWhere, 0, len(aFilters))

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
		oWhere := &pkgClickhouse.ClickhouseWhere{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oValue,
		}
		aWheres = append(aWheres, oWhere)
	}

	return aWheres
}

func (oSelf *AbstractClickhouse) SortersToOrders(aSorters []*pkgInput.Sorter) []*pkgClickhouse.ClickhouseOrder {
	aOrders := make([]*pkgClickhouse.ClickhouseOrder, 0, len(aSorters))

	for _, oSorter := range aSorters {
		if oSorter == nil || oSorter.Field == nil {
			continue
		}

		sDirection := "asc"
		if oSorter.Order != nil && strings.EqualFold(*oSorter.Order, "desc") {
			sDirection = "desc"
		}

		sField := "`" + *oSorter.Field + "`"
		oOrder := &pkgClickhouse.ClickhouseOrder{
			Field: &sField,
			Value: &sDirection,
		}
		aOrders = append(aOrders, oOrder)
	}

	return aOrders
}

func (oSelf *AbstractClickhouse) PaginationToLimit(oPagination *pkgInput.Pagination) *pkgClickhouse.ClickhouseLimit {
	iSize := uint(10)
	iPage := uint(1)

	if oPagination != nil && oPagination.Size != nil && *oPagination.Size != 0 {
		iSize = *oPagination.Size
	}

	if oPagination != nil && oPagination.Page != nil && *oPagination.Page != 0 {
		iPage = *oPagination.Page
	}

	iOffset := (iPage - 1) * iSize

	return &pkgClickhouse.ClickhouseLimit{
		Offset: &iOffset,
		Count:  &iSize,
	}
}
