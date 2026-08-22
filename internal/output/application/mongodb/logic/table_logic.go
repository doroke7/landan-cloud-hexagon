package logic

import (
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"

	domain "example/internal/domain"
	outputApplicationMongodb "example/internal/output/application/mongodb"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkg "example/pkg"
)

type TableLogic struct {
	*outputApplicationMongodb.AbstractMongodb
	Collection *mongo.Collection
}

func NewTableLogic(oAbstractMongodb *outputApplicationMongodb.AbstractMongodb) outputPortAnyLogic.TableLogic {
	return &TableLogic{
		AbstractMongodb: oAbstractMongodb,
		Collection:      oAbstractMongodb.Database.Collection("tables"),
	}
}

func (oSelf *TableLogic) ShowTablesTotalByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Table, int64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersPaginationToFindOptions(aSorters, oPagination)

	var aDocs []*domain.TableDocument
	var iTotal int64
	var oFindErr error
	var oCountErr error

	var oWaitGroup sync.WaitGroup
	oWaitGroup.Add(2)

	go func() {
		defer oWaitGroup.Done()

		oCursor, oErr := oSelf.Collection.Find(oSelf.Context, oFilter, oFindOptions)
		if oErr != nil {
			oFindErr = oErr
			return
		}
		defer oCursor.Close(oSelf.Context)

		oFindErr = oCursor.All(oSelf.Context, &aDocs)
	}()

	go func() {
		defer oWaitGroup.Done()

		iTotal, oCountErr = oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	}()

	oWaitGroup.Wait()

	aTables := make([]*domain.Table, len(aDocs))
	for i, oDoc := range aDocs {
		aTables[i] = domain.TableDocumentToTable(oDoc)
	}

	if oFindErr != nil {
		return aTables, 0, oFindErr
	}

	if oCountErr != nil {
		return aTables, 0, oCountErr
	}

	return aTables, iTotal, nil
}
