package outputApplicationGaussdb

import (
	"context"
	"strings"

	pkgCache "example/pkg/cache"
	pkgGaussdb "example/pkg/gaussdb"
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
type AbstractGaussdb struct {
	DB      *gorm.DB
	Context context.Context
	*pkgCache.Aop
}

func NewAbstractGaussdb(oContext context.Context, oDb *gorm.DB, oAop *pkgCache.Aop) *AbstractGaussdb {

	return &AbstractGaussdb{
		DB:      oDb,
		Context: oContext,
		Aop:     oAop,
	}
}

func (oSelf *AbstractGaussdb) FiltersToWheres(aFilters []*pkgInput.Filter) []*pkgGaussdb.GaussdbWhere {
	aWheres := make([]*pkgGaussdb.GaussdbWhere, 0, len(aFilters))

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

		// GaussDB 源自 PostgreSQL，識別字一樣用雙引號。
		sField := `"` + *oFilter.Field + `"`
		oWhere := &pkgGaussdb.GaussdbWhere{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oValue,
		}
		aWheres = append(aWheres, oWhere)
	}

	return aWheres
}

func (oSelf *AbstractGaussdb) SortersToOrders(aSorters []*pkgInput.Sorter) []*pkgGaussdb.GaussdbOrder {
	aOrders := make([]*pkgGaussdb.GaussdbOrder, 0, len(aSorters))

	for _, oSorter := range aSorters {
		if oSorter == nil || oSorter.Field == nil {
			continue
		}

		sDirection := "asc"
		if oSorter.Order != nil && strings.EqualFold(*oSorter.Order, "desc") {
			sDirection = "desc"
		}

		sField := `"` + *oSorter.Field + `"`
		oOrder := &pkgGaussdb.GaussdbOrder{
			Field: &sField,
			Value: &sDirection,
		}
		aOrders = append(aOrders, oOrder)
	}

	return aOrders
}

func (oSelf *AbstractGaussdb) PaginationToLimit(oPagination *pkgInput.Pagination) *pkgGaussdb.GaussdbLimit {
	iSize := uint(10)
	iPage := uint(1)

	if oPagination != nil && oPagination.Size != nil && *oPagination.Size != 0 {
		iSize = *oPagination.Size
	}

	if oPagination != nil && oPagination.Page != nil && *oPagination.Page != 0 {
		iPage = *oPagination.Page
	}

	iOffset := (iPage - 1) * iSize

	return &pkgGaussdb.GaussdbLimit{
		Offset: &iOffset,
		Count:  &iSize,
	}
}
