package mysql

import (
	"context"

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
type AbstractModel struct {
	DB      *gorm.DB
	Context context.Context
	*pkg.Aop
}

func NewAbstractModel(oContext context.Context, oDb *gorm.DB, oAop *pkg.Aop) *AbstractModel {

	return &AbstractModel{
		DB:      oDb,
		Context: oContext,
		Aop:     oAop,
	}
}

func (oSelf *AbstractModel) FiltersToWheres(aFields []string, aFilters []*pkg.Filter) []*pkg.Where {
	bAllowAll := len(aFields) == 0

	oAllowed := make(map[string]bool, len(aFields))
	for _, sField := range aFields {
		oAllowed[sField] = true
	}

	aWheres := make([]*pkg.Where, 0, len(aFilters))

	for _, oFilter := range aFilters {
		if oFilter == nil || oFilter.Field == nil || (!bAllowAll && !oAllowed[*oFilter.Field]) {
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

		sField := "`" + *oFilter.Field + "`"
		aWheres = append(aWheres, &pkg.Where{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oValue,
		})
	}

	return aWheres
}
