package outputApplicationPostgresql

import (
	"context"
	"strings"

	pkgCache "example/pkg/cache"
	pkgInput "example/pkg/input"
	pkgPostgresql "example/pkg/postgresql"

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
type AbstractPostgresql struct {
	DB      *gorm.DB
	Context context.Context
	*pkgCache.Aop
}

func NewAbstractPostgresql(oContext context.Context, oDb *gorm.DB, oAop *pkgCache.Aop) *AbstractPostgresql {

	return &AbstractPostgresql{
		DB:      oDb,
		Context: oContext,
		Aop:     oAop,
	}
}

func (oSelf *AbstractPostgresql) FiltersToWheres(aFilters []*pkgInput.Filter) []*pkgPostgresql.PostgresqlWhere {
	aWheres := make([]*pkgPostgresql.PostgresqlWhere, 0, len(aFilters))

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

		// Postgres 的識別字用雙引號，跟 mysql 的反引號不同。
		sField := `"` + *oFilter.Field + `"`
		aWheres = append(aWheres, &pkgPostgresql.PostgresqlWhere{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oValue,
		})
	}

	return aWheres
}

func (oSelf *AbstractPostgresql) SortersToOrders(aSorters []*pkgInput.Sorter) []*pkgPostgresql.PostgresqlOrder {
	aOrders := make([]*pkgPostgresql.PostgresqlOrder, 0, len(aSorters))

	for _, oSorter := range aSorters {
		if oSorter == nil || oSorter.Field == nil {
			continue
		}

		sDirection := "asc"
		if oSorter.Order != nil && strings.EqualFold(*oSorter.Order, "desc") {
			sDirection = "desc"
		}

		sField := `"` + *oSorter.Field + `"`
		aOrders = append(aOrders, &pkgPostgresql.PostgresqlOrder{
			Field: &sField,
			Value: &sDirection,
		})
	}

	return aOrders
}

func (oSelf *AbstractPostgresql) PaginationToLimit(oPagination *pkgInput.Pagination) *pkgPostgresql.PostgresqlLimit {
	iSize := uint(10)
	iPage := uint(1)

	if oPagination != nil && oPagination.Size != nil && *oPagination.Size != 0 {
		iSize = *oPagination.Size
	}

	if oPagination != nil && oPagination.Page != nil && *oPagination.Page != 0 {
		iPage = *oPagination.Page
	}

	iOffset := (iPage - 1) * iSize

	return &pkgPostgresql.PostgresqlLimit{
		Offset: &iOffset,
		Count:  &iSize,
	}
}
