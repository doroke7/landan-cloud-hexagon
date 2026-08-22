package pkg

import (
	"fmt"
	"strings"
)

func SortersToMysqlOrders(aFields []string, aSorters []*Sorter) []*MysqlOrder {
	bAllowAll := len(aFields) == 0

	oAllowed := make(map[string]bool, len(aFields))
	for _, sField := range aFields {
		oAllowed[sField] = true
	}

	aOrders := make([]*MysqlOrder, 0, len(aSorters))

	for _, oSorter := range aSorters {
		fmt.Printf("%#v", oSorter)
		fmt.Println("oSorter=", oSorter)

		if oSorter == nil || oSorter.Field == nil || (!bAllowAll && !oAllowed[*oSorter.Field]) {
			continue
		}

		sDirection := "asc"
		if oSorter.Order != nil && strings.EqualFold(*oSorter.Order, "desc") {
			sDirection = "desc"
		}

		sField := "`" + *oSorter.Field + "`"
		aOrders = append(aOrders, &MysqlOrder{
			Field: &sField,
			Value: &sDirection,
		})
	}

	return aOrders
}
