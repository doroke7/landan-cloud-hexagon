package resource

import (
	"context"

	client "example/internal/client"
	pbResource "example/pb/resource"
	pkg "example/pkg"

	"google.golang.org/protobuf/types/known/structpb"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractRepository 做法一致。
type AbstractLogic struct {
	Context             context.Context
	ResourceLogicClient *client.Logic
}

func NewAbstractLogic(oContext context.Context, oResourceClient *client.ResourceClient) *AbstractLogic {
	return &AbstractLogic{
		Context:             oContext,
		ResourceLogicClient: oResourceClient.Logic,
	}
}

func (oSelf *AbstractLogic) ToFilters(aFilters []*pkg.Filter) []*pbResource.Filter {
	aPbFilters := make([]*pbResource.Filter, 0, len(aFilters))

	for _, oFilter := range aFilters {
		if oFilter == nil || oFilter.Field == nil {
			continue
		}

		oValue, oErr := structpb.NewValue(oFilter.Value)
		if oErr != nil {
			continue
		}

		oPbFilter := &pbResource.Filter{
			Field: *oFilter.Field,
			Value: oValue,
		}
		if oFilter.Operator != nil {
			oPbFilter.Operator = *oFilter.Operator
		}

		aPbFilters = append(aPbFilters, oPbFilter)
	}

	return aPbFilters
}

func (oSelf *AbstractLogic) ToSorters(aSorters []*pkg.Sorter) []*pbResource.Sorter {
	aPbSorters := make([]*pbResource.Sorter, 0, len(aSorters))

	for _, oSorter := range aSorters {
		if oSorter == nil || oSorter.Field == nil || oSorter.Order == nil {
			continue
		}

		aPbSorters = append(aPbSorters, &pbResource.Sorter{
			Field: *oSorter.Field,
			Order: *oSorter.Order,
		})
	}

	return aPbSorters
}

func (oSelf *AbstractLogic) ToPagination(oPagination *pkg.Pagination) *pbResource.Pagination {
	if oPagination == nil {
		return nil
	}

	oPbPagination := &pbResource.Pagination{}

	if oPagination.Size != nil {
		oPbPagination.Size = uint64(*oPagination.Size)
	}
	if oPagination.Page != nil {
		oPbPagination.Page = uint64(*oPagination.Page)
	}

	return oPbPagination
}
