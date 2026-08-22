package pkg

// eq, ne, gt, gte, lt, lte, contains, startsWith, endsWith, notContains, in, notIn, between
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

func FiltersToWheres(aFields []string, aFilters []*Filter) []*Where {
	bAllowAll := len(aFields) == 0

	oAllowed := make(map[string]bool, len(aFields))
	for _, sField := range aFields {
		oAllowed[sField] = true
	}

	aWheres := make([]*Where, 0, len(aFilters))

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
		aWheres = append(aWheres, &Where{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oValue,
		})
	}

	return aWheres
}
