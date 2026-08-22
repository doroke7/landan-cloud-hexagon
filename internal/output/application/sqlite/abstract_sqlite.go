package sqlite

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
type AbstractSqlite struct {
	DB      *gorm.DB
	Context context.Context
	*pkg.Aop
}

func NewAbstractSqlite(oContext context.Context, oDb *gorm.DB, oAop *pkg.Aop) *AbstractSqlite {

	return &AbstractSqlite{
		DB:      oDb,
		Context: oContext,
		Aop:     oAop,
	}
}

func (oSelf *AbstractSqlite) FiltersToWheres(aFilters []*pkg.Filter) []*pkg.SqliteWhere {
	aWheres := make([]*pkg.SqliteWhere, 0, len(aFilters))

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

		// SQLite 的識別字用雙引號，跟 postgres 一樣是 ANSI 標準寫法。
		sField := `"` + *oFilter.Field + `"`
		aWheres = append(aWheres, &pkg.SqliteWhere{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oValue,
		})
	}

	return aWheres
}

func (oSelf *AbstractSqlite) SortersToOrders(aSorters []*pkg.Sorter) []*pkg.SqliteOrder {
	aOrders := make([]*pkg.SqliteOrder, 0, len(aSorters))

	for _, oSorter := range aSorters {
		if oSorter == nil || oSorter.Field == nil {
			continue
		}

		sDirection := "asc"
		if oSorter.Order != nil && strings.EqualFold(*oSorter.Order, "desc") {
			sDirection = "desc"
		}

		sField := `"` + *oSorter.Field + `"`
		aOrders = append(aOrders, &pkg.SqliteOrder{
			Field: &sField,
			Value: &sDirection,
		})
	}

	return aOrders
}

func (oSelf *AbstractSqlite) PaginationToLimit(oPagination *pkg.Pagination) *pkg.SqliteLimit {
	iSize := uint(10)
	iPage := uint(1)

	if oPagination != nil && oPagination.Size != nil && *oPagination.Size != 0 {
		iSize = *oPagination.Size
	}

	if oPagination != nil && oPagination.Page != nil && *oPagination.Page != 0 {
		iPage = *oPagination.Page
	}

	iOffset := (iPage - 1) * iSize

	return &pkg.SqliteLimit{
		Offset: &iOffset,
		Count:  &iSize,
	}
}
