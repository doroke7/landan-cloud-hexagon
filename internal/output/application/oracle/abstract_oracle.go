package outputApplicationOracle

import (
	"context"
	"strings"

	pkgCache "example/pkg/cache"
	pkgInput "example/pkg/input"
	pkgOracle "example/pkg/oracle"

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
type AbstractOracle struct {
	DB      *gorm.DB
	Context context.Context
	*pkgCache.Aop
}

func NewAbstractOracle(oContext context.Context, oDb *gorm.DB, oAop *pkgCache.Aop) *AbstractOracle {

	return &AbstractOracle{
		DB:      oDb,
		Context: oContext,
		Aop:     oAop,
	}
}

func (oSelf *AbstractOracle) FiltersToWheres(aFilters []*pkgInput.Filter) []*pkgOracle.OracleWhere {
	aWheres := make([]*pkgOracle.OracleWhere, 0, len(aFilters))

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

		// Oracle 的識別字用雙引號，跟 postgres 一樣是 ANSI 標準寫法。
		sField := `"` + *oFilter.Field + `"`
		aWheres = append(aWheres, &pkgOracle.OracleWhere{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oValue,
		})
	}

	return aWheres
}

func (oSelf *AbstractOracle) SortersToOrders(aSorters []*pkgInput.Sorter) []*pkgOracle.OracleOrder {
	aOrders := make([]*pkgOracle.OracleOrder, 0, len(aSorters))

	for _, oSorter := range aSorters {
		if oSorter == nil || oSorter.Field == nil {
			continue
		}

		sDirection := "asc"
		if oSorter.Order != nil && strings.EqualFold(*oSorter.Order, "desc") {
			sDirection = "desc"
		}

		sField := `"` + *oSorter.Field + `"`
		aOrders = append(aOrders, &pkgOracle.OracleOrder{
			Field: &sField,
			Value: &sDirection,
		})
	}

	return aOrders
}

func (oSelf *AbstractOracle) PaginationToLimit(oPagination *pkgInput.Pagination) *pkgOracle.OracleLimit {
	iSize := uint(10)
	iPage := uint(1)

	if oPagination != nil && oPagination.Size != nil && *oPagination.Size != 0 {
		iSize = *oPagination.Size
	}

	if oPagination != nil && oPagination.Page != nil && *oPagination.Page != 0 {
		iPage = *oPagination.Page
	}

	iOffset := (iPage - 1) * iSize

	return &pkgOracle.OracleLimit{
		Offset: &iOffset,
		Count:  &iSize,
	}
}
