package outputApplicationMongodbLogic

import (
	"errors"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupLogic struct {
	*AbstractLogic
	Collection       *mongo.Collection
	AdminPermissions *mongo.Collection
	Counters         *mongo.Collection
}

func NewAdminPermissionGroupLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminPermissionGroupLogic {
	oCollection := oAbstractLogic.Database.Collection("admin_permission_groups")
	oAdminPermissions := oAbstractLogic.Database.Collection("admin_permissions")
	oCounters := oAbstractLogic.Database.Collection("counters")

	return &AdminPermissionGroupLogic{
		AbstractLogic:    oAbstractLogic,
		Collection:       oCollection,
		AdminPermissions: oAdminPermissions,
		Counters:         oCounters,
	}
}

func (oSelf *AdminPermissionGroupLogic) nextId() (uint64, error) {
	oUpdateOptions := options.FindOneAndUpdate()
	oUpdateOptions.SetUpsert(true)
	oUpdateOptions.SetReturnDocument(options.After)

	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "admin_permission_group"},
		bson.M{"$inc": bson.M{"seq": 1}},
		oUpdateOptions,
	)

	var oCounter struct {
		Seq uint64 `bson:"seq"`
	}
	if oErr := oResult.Decode(&oCounter); oErr != nil {
		return 0, oErr
	}

	return oCounter.Seq, nil
}

func (oSelf *AdminPermissionGroupLogic) nextAdminPermissionId() (uint64, error) {
	oUpdateOptions := options.FindOneAndUpdate()
	oUpdateOptions.SetUpsert(true)
	oUpdateOptions.SetReturnDocument(options.After)

	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "admin_permission"},
		bson.M{"$inc": bson.M{"seq": 1}},
		oUpdateOptions,
	)

	var oCounter struct {
		Seq uint64 `bson:"seq"`
	}
	if oErr := oResult.Decode(&oCounter); oErr != nil {
		return 0, oErr
	}

	return oCounter.Seq, nil
}

func (oSelf *AdminPermissionGroupLogic) AddAdminPermissionGroup(oValue *domain.AdminPermissionGroupVariable) error {
	iId, oErr := oSelf.nextId()
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oNew := &domain.AdminPermissionGroup{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oValue.Key != nil {
		oNew.Key = *oValue.Key
	}
	if oValue.Name != nil {
		oNew.Name = *oValue.Name
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oNew); oErr != nil {
		return oErr
	}

	for _, oAdminPermissionValue := range oValue.AdminPermissions {
		if oAdminPermissionValue == nil {
			continue
		}

		if oErr := oSelf.insertAdminPermission(oAdminPermissionValue, iId); oErr != nil {
			return oErr
		}
	}

	return nil
}

func (oSelf *AdminPermissionGroupLogic) insertAdminPermission(oValue *domain.AdminPermissionVariable, iAdminPermissionGroupId uint64) error {
	iId, oErr := oSelf.nextAdminPermissionId()
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oNew := &domain.AdminPermission{
		Id:                     iId,
		AdminPermissionGroupId: iAdminPermissionGroupId,
		CreatedAt:              oNow,
		UpdatedAt:              oNow,
		DeletedAt:              oDeletedAtZero,
	}

	if oValue.Type != nil {
		oNew.Type = *oValue.Type
	}
	if oValue.Key != nil {
		oNew.Key = *oValue.Key
	}
	if oValue.Name != nil {
		oNew.Name = *oValue.Name
	}

	_, oErr = oSelf.AdminPermissions.InsertOne(oSelf.Context, oNew)

	return oErr
}

// EditAdminPermissionGroupById 更新 group 本身欄位，並同步 admin_permissions：
// 傳進來沒帶 id 的就新增，帶 id 的就修改；db 裡面現有、但沒出現在傳進來 id 清單裡的就刪除。
func (oSelf *AdminPermissionGroupLogic) EditAdminPermissionGroupById(oValue *domain.AdminPermissionGroupVariable, iId uint64) error {
	oSet := bson.M{"updated_at": time.Now()}

	if oValue.Key != nil {
		oSet["key"] = *oValue.Key
	}
	if oValue.Name != nil {
		oSet["name"] = *oValue.Name
	}

	oResult, oErr := oSelf.Collection.UpdateOne(
		oSelf.Context,
		bson.M{"_id": iId},
		bson.M{"$set": oSet},
	)
	if oErr != nil {
		return oErr
	}

	if oResult.MatchedCount == 0 {
		return errors.New("0 rows updated")
	}

	// 撈出 db 現有未刪除的 admin_permissions，準備跟傳進來的 id 清單比對
	oCursor, oErr := oSelf.AdminPermissions.Find(oSelf.Context, bson.M{
		"admin_permission_group_id": iId,
		"deleted_at":                oDeletedAtZero,
	})
	if oErr != nil {
		return oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aExistingAdminPermissions []*domain.AdminPermission
	if oErr := oCursor.All(oSelf.Context, &aExistingAdminPermissions); oErr != nil {
		return oErr
	}

	// 逐筆處理傳進來的 admin_permissions：沒帶 id 就新增，帶 id 就修改
	aKeptAdminPermissionIds := make(map[uint64]bool, len(oValue.AdminPermissions))
	for _, oAdminPermissionValue := range oValue.AdminPermissions {
		if oAdminPermissionValue == nil {
			continue
		}

		if oAdminPermissionValue.Id == nil {
			if oErr := oSelf.insertAdminPermission(oAdminPermissionValue, iId); oErr != nil {
				return oErr
			}

			continue
		}

		iAdminPermissionId := *oAdminPermissionValue.Id
		aKeptAdminPermissionIds[iAdminPermissionId] = true

		oAdminPermissionSet := bson.M{"updated_at": time.Now()}
		if oAdminPermissionValue.Type != nil {
			oAdminPermissionSet["type"] = *oAdminPermissionValue.Type
		}
		if oAdminPermissionValue.Key != nil {
			oAdminPermissionSet["key"] = *oAdminPermissionValue.Key
		}
		if oAdminPermissionValue.Name != nil {
			oAdminPermissionSet["name"] = *oAdminPermissionValue.Name
		}

		_, oErr := oSelf.AdminPermissions.UpdateOne(
			oSelf.Context,
			bson.M{"_id": iAdminPermissionId},
			bson.M{"$set": oAdminPermissionSet},
		)
		if oErr != nil {
			return oErr
		}
	}

	// db 裡面有、但沒出現在傳進來 id 清單裡的，軟刪除
	oNow := time.Now()
	for _, oExistingAdminPermission := range aExistingAdminPermissions {
		if aKeptAdminPermissionIds[oExistingAdminPermission.Id] {
			continue
		}

		_, oErr := oSelf.AdminPermissions.UpdateOne(
			oSelf.Context,
			bson.M{"_id": oExistingAdminPermission.Id},
			bson.M{"$set": bson.M{"deleted_at": oNow}},
		)
		if oErr != nil {
			return oErr
		}
	}

	return nil
}

// ShowTree 先把所有未刪除的 admin_permission_group 一次撈成平的，再用 ParentId 掛 Children，
// 回傳 ParentId == 0 的 root。
func (oSelf *AdminPermissionGroupLogic) ShowTree() ([]*domain.AdminPermissionGroup, error) {
	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, bson.M{"deleted_at": oDeletedAtZero})
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aFlat []*domain.AdminPermissionGroup
	if oErr := oCursor.All(oSelf.Context, &aFlat); oErr != nil {
		return nil, oErr
	}

	aByParent := make(map[uint64][]*domain.AdminPermissionGroup, len(aFlat))
	for _, oOne := range aFlat {
		aByParent[oOne.ParentId] = append(aByParent[oOne.ParentId], oOne)
	}

	var fnAttach func(oNode *domain.AdminPermissionGroup)
	fnAttach = func(oNode *domain.AdminPermissionGroup) {
		for _, oChild := range aByParent[oNode.Id] {
			fnAttach(oChild)
			oNode.Children = append(oNode.Children, oChild)
		}
	}

	aRoots := aByParent[0]
	for _, oRoot := range aRoots {
		fnAttach(oRoot)
	}

	return aRoots, nil
}

func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroupById(iId uint64) (*domain.AdminPermissionGroup, error) {
	var oAdminPermissionGroup domain.AdminPermissionGroup

	oResult := oSelf.Collection.FindOne(oSelf.Context, bson.M{
		"_id":        iId,
		"deleted_at": oDeletedAtZero,
	})
	oErr := oResult.Decode(&oAdminPermissionGroup)

	if oErr != nil {
		if oErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, oErr
	}

	return &oAdminPermissionGroup, nil
}

func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroups() ([]*domain.AdminPermissionGroup, error) {
	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, bson.M{"deleted_at": oDeletedAtZero})
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aAdminPermissionGroups []*domain.AdminPermissionGroup
	if oErr := oCursor.All(oSelf.Context, &aAdminPermissionGroups); oErr != nil {
		return nil, oErr
	}

	return aAdminPermissionGroups, nil
}

func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersPaginationToFindOptions(aSorters, oPagination)

	var aAdminPermissionGroups []*domain.AdminPermissionGroup
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

		oFindErr = oCursor.All(oSelf.Context, &aAdminPermissionGroups)
	}()

	go func() {
		defer oWaitGroup.Done()

		iTotal, oCountErr = oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	}()

	oWaitGroup.Wait()

	if oFindErr != nil {
		return aAdminPermissionGroups, 0, oFindErr
	}

	if oCountErr != nil {
		return aAdminPermissionGroups, 0, oCountErr
	}

	return aAdminPermissionGroups, uint64(iTotal), nil
}
