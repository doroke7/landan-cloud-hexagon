package pkg

import (
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func SortersToSort(aSorters []*Sorter) bson.D {
	oSort := bson.D{}

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

	return oSort
}
