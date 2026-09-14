package outputApplicationResource

import (
	pbResource "example/pb/resource"
	pkgInput "example/pkg/input"

	"google.golang.org/protobuf/types/known/structpb"
)

// AbstractResource 只放跟 resource client 無關的純轉換方法（pkgInput -> pb），
// 由各層的 AbstractModel / AbstractLogic / AbstractEvent 內嵌後方法提升上去；
// client 欄位留在各層自己的 abstract 裡。
type AbstractResource struct{}

func NewAbstractResource() *AbstractResource {
	oResource := &AbstractResource{}

	return oResource
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

		oPbSorter := &pbResource.Sorter{
			Field: *oSorter.Field,
			Order: *oSorter.Order,
		}
		aPbSorters = append(aPbSorters, oPbSorter)
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
