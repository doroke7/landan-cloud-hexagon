package mongodb

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	bootstrap "example/bootstrap"
	pkg "example/pkg"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractModel 做法一致。
// Database 是從 bootstrap 注入的共用 mongo 連線取出的資料庫控制代碼（CONFIG.MONGODB.NAME），
// 跟 mysql.AbstractMysql 持有 *gorm.DB 是同一種角色。
type AbstractMongodb struct {
	Context  context.Context
	Client   *mongo.Client
	Database *mongo.Database
}

func NewAbstractMongodb(oContext context.Context, oClient *mongo.Client) *AbstractMongodb {
	return &AbstractMongodb{
		Context:  oContext,
		Client:   oClient,
		Database: oClient.Database(bootstrap.CONFIG.MONGODB.NAME),
	}
}

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

func (oSelf *AbstractMongodb) FiltersToFilter(aFilters []*pkg.Filter) bson.M {
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

func (oSelf *AbstractMongodb) SortersToFindOptions(aSorters []*pkg.Sorter) *options.FindOptionsBuilder {
	oSort := bson.D{}
	oFindOptions := options.Find()

	for _, oSorter := range aSorters {
		if oSorter == nil || oSorter.Field == nil {
			continue
		}

		iDirection := 1
		if oSorter.Order != nil && strings.EqualFold(*oSorter.Order, "desc") {
			iDirection = -1
		}

		oSort = append(oSort, bson.E{Key: *oSorter.Field, Value: iDirection})
	}

	if len(oSort) > 0 {
		oFindOptions.SetSort(oSort)
	}

	return oFindOptions
}

func (oSelf *AbstractMongodb) SortersPaginationToFindOptions(aSorters []*pkg.Sorter, oPagination *pkg.Pagination) *options.FindOptionsBuilder {
	oSort := bson.D{}
	oFindOptions := options.Find()

	for _, oSorter := range aSorters {
		if oSorter == nil || oSorter.Field == nil {
			continue
		}

		iDirection := 1
		if oSorter.Order != nil && strings.EqualFold(*oSorter.Order, "desc") {
			iDirection = -1
		}

		oSort = append(oSort, bson.E{Key: *oSorter.Field, Value: iDirection})
	}

	if len(oSort) > 0 {
		oFindOptions.SetSort(oSort)
	}

	iSize := uint(10)
	iPage := uint(1)
	if oPagination != nil {
		if oPagination.Size != nil && *oPagination.Size != 0 {
			iSize = *oPagination.Size
		}
		if oPagination.Page != nil && *oPagination.Page != 0 {
			iPage = *oPagination.Page
		}
	}
	oFindOptions.SetLimit(int64(iSize)).SetSkip(int64((iPage - 1) * iSize))

	return oFindOptions
}
