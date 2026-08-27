package outputApplicationResource

import (
	"context"

	client "example/internal/client"
	pbResource "example/pb/resource"
	pkgInput "example/pkg/input"

	"google.golang.org/protobuf/types/known/structpb"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractRepository 做法一致。
type AbstractResource struct {
	Context             context.Context
	ResourceLogicClient *client.Logic
	ResourceModelClient *client.Model
}

func NewAbstractResource(oContext context.Context, oResourceClient *client.ResourceClient) *AbstractResource {
	return &AbstractResource{
		Context:             oContext,
		ResourceLogicClient: oResourceClient.Logic,
		ResourceModelClient: oResourceClient.Model,
	}
}

func (oSelf *AbstractResource) ToFilters(aFilters []*pkgInput.Filter) []*pbResource.Filter {
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

func (oSelf *AbstractResource) ToSorters(aSorters []*pkgInput.Sorter) []*pbResource.Sorter {
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

func (oSelf *AbstractResource) ToPagination(oPagination *pkgInput.Pagination) *pbResource.Pagination {
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
