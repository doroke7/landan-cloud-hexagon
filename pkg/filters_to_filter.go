package pkg

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// eq, ne, gt, gte, lt, lte, in, notIn 直接對應 Mongo 的比較運算子，
// contains, notContains, startsWith, endsWith 用不分大小寫的 regex 模擬，
// between 預期 Value 是長度為 2 的 slice，轉成 $gte/$lte 區間。
var oFilterOperatorMap = map[string]string{
	"eq":    "$eq",
	"ne":    "$ne",
	"gt":    "$gt",
	"gte":   "$gte",
	"lt":    "$lt",
	"lte":   "$lte",
	"in":    "$in",
	"notIn": "$nin",
}

func FiltersToFilter(aFilters []*Filter) bson.M {
	oFilter := bson.M{}

	for _, oFilter1 := range aFilters {
		if oFilter1 == nil || oFilter1.Field == nil {
			continue
		}

		sField := *oFilter1.Field

		if oFilter1.Operator == nil {
			oFilter[sField] = oFilter1.Value
			continue
		}

		switch *oFilter1.Operator {
		case "contains":
			if sValue, bOk := oFilter1.Value.(string); bOk {
				oFilter[sField] = bson.M{"$regex": sValue, "$options": "i"}
			}
		case "notContains":
			if sValue, bOk := oFilter1.Value.(string); bOk {
				oFilter[sField] = bson.M{"$not": bson.M{"$regex": sValue, "$options": "i"}}
			}
		case "startsWith":
			if sValue, bOk := oFilter1.Value.(string); bOk {
				oFilter[sField] = bson.M{"$regex": "^" + sValue, "$options": "i"}
			}
		case "endsWith":
			if sValue, bOk := oFilter1.Value.(string); bOk {
				oFilter[sField] = bson.M{"$regex": sValue + "$", "$options": "i"}
			}
		case "between":
			if aRange, bOk := oFilter1.Value.([]any); bOk && len(aRange) == 2 {
				oFilter[sField] = bson.M{"$gte": aRange[0], "$lte": aRange[1]}
			}
		default:
			if sOperator, bOk := oFilterOperatorMap[*oFilter1.Operator]; bOk {
				oFilter[sField] = bson.M{sOperator: oFilter1.Value}
			} else {
				oFilter[sField] = oFilter1.Value
			}
		}
	}

	return oFilter
}
