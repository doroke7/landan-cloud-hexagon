package outputApplicationElasticsearchLogic

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type AdminUserLogic struct {
	*AbstractLogic
	Index string
}

func NewAdminUserLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminUserLogic {
	return &AdminUserLogic{
		AbstractLogic: oAbstractLogic,
		Index:         oAbstractLogic.IndexName("admin_users"),
	}
}

// ShowAdminUsersTotalByFiltersWithSortersPagination 一次 search 就同時拿到「這一頁」跟「符合條件的總數」
// （options 裡有 track_total_hits），不用再另外打一次 count。
func (oSelf *AdminUserLogic) ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint64, error) {
	sDeletedAtField := "deleted_at"
	aFilters = append(aFilters, &pkgInput.Filter{Field: &sDeletedAtField, Value: oDeletedAtZero})

	aOptions, oErr := oSelf.IndexFiltersSortersPaginationToOptions(oSelf.Index, aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, 0, oErr
	}

	oResult, oErr := oSelf.SearchWithOptions(aOptions)
	if oErr != nil {
		return nil, 0, oErr
	}

	aAdminUsers := make([]*domain.AdminUser, len(oResult.Hits))
	for i, oHit := range oResult.Hits {
		var oAdminUser domain.AdminUser
		if oErr := json.Unmarshal(oHit.Source, &oAdminUser); oErr != nil {
			return nil, 0, oErr
		}
		aAdminUsers[i] = &oAdminUser
	}

	return aAdminUsers, uint64(oResult.Total), nil
}

func (oSelf *AdminUserLogic) AddAminUser(oValue *domain.AdminUserVariable) error {
	iId, oErr := oSelf.NextId("admin_user")
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oDoc := &domain.AdminUser{
		Id:        uint64(iId),
		CreatedAt: oNow,
		UpdatedAt: oNow,
	}

	if oValue.Name != nil {
		oDoc.Name = *oValue.Name
	}
	if oValue.Password != nil {
		oDoc.Password = *oValue.Password
	}

	if oErr := oSelf.IndexOne(oSelf.Index, strconv.FormatUint(uint64(iId), 10), oDoc); oErr != nil {
		return oErr
	}

	return nil
}

func (oSelf *AdminUserLogic) EditAdminUserById(oValue *domain.AdminUserVariable, iId uint64) error {
	oColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}
	delete(oColumns, "admin_role_ids")
	oColumns["updated_at"] = time.Now()

	sId := strconv.FormatUint(uint64(iId), 10)

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oColumns)
	if oErr != nil {
		return oErr
	}

	if !bOk {
		return errors.New("0 rows updated")
	}

	return nil
}

func (oSelf *AdminUserLogic) ShowAdminUserById(iId uint64) (*domain.AdminUser, error) {
	var oAdminUser domain.AdminUser

	sId := strconv.FormatUint(iId, 10)

	bFound, oErr := oSelf.GetById(oSelf.Index, sId, &oAdminUser)
	if oErr != nil {
		return nil, oErr
	}

	if !bFound {
		return nil, errors.New("record not found")
	}

	return &oAdminUser, nil
}

func (oSelf *AdminUserLogic) RemoveAdminUserById(iId uint64) error {
	sId := strconv.FormatUint(iId, 10)
	oPartial := map[string]any{"deleted_at": time.Now()}

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oPartial)
	if oErr != nil {
		return oErr
	}

	if !bOk {
		return errors.New("0 rows deleted")
	}

	return nil
}
