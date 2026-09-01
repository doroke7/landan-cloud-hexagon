package outputApplicationElasticsearchModel

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"

	domain "example/internal/domain"
	elasticsearchBase "example/internal/output/application/elasticsearch"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AppUserModel struct {
	*elasticsearchBase.AbstractElasticsearch
	Index string
}

func NewAppUserModel(oAbstractModel *elasticsearchBase.AbstractElasticsearch) outputPortAnyModel.AppUserModel {
	return &AppUserModel{
		AbstractElasticsearch: oAbstractModel,
		Index:                 oAbstractModel.IndexName("app_users"),
	}
}

// IncreaseBalance 讀出目前餘額、加上 amount 再寫回。ES 沒有 gorm 那種 UpdateColumn 原子加，
// 這裡是「讀-改-寫」，不保證併發原子性（跟 lottery/memory adapter 的取捨一致）。
func (oSelf *AppUserModel) IncreaseBalance(iId uint, iAmount uint) (bool, error) {
	oAppUser, oErr := oSelf.ShowOneById(iId)
	if oErr != nil {
		return false, oErr
	}
	if oAppUser == nil {
		return false, errors.New("資料不存在")
	}

	sId := strconv.FormatUint(uint64(iId), 10)
	oPartial := map[string]any{
		"balance":    oAppUser.Balance + iAmount,
		"updated_at": time.Now(),
	}

	return oSelf.UpdateOne(oSelf.Index, sId, oPartial)
}

func (oSelf *AppUserModel) ShowOneByName(sName string) (*domain.AppUser, error) {
	oBody := map[string]any{
		"query": map[string]any{
			"term": map[string]any{"name": sName},
		},
	}

	aBodyBytes, oErr := json.Marshal(oBody)
	if oErr != nil {
		return nil, oErr
	}

	aOptions := []func(*esapi.SearchRequest){
		oSelf.Client.Search.WithContext(oSelf.Context),
		oSelf.Client.Search.WithIndex(oSelf.Index),
		oSelf.Client.Search.WithBody(bytes.NewReader(aBodyBytes)),
		oSelf.Client.Search.WithTrackTotalHits(true),
		oSelf.Client.Search.WithSize(1),
	}

	oResult, oErr := oSelf.SearchWithOptions(aOptions)
	if oErr != nil {
		return nil, oErr
	}

	if len(oResult.Hits) == 0 {
		return nil, errors.New("資料不存在")
	}

	var oAppUser domain.AppUser
	if oErr := json.Unmarshal(oResult.Hits[0].Source, &oAppUser); oErr != nil {
		return nil, oErr
	}

	return &oAppUser, nil
}

func (oSelf *AppUserModel) ShowOneById(iId uint) (*domain.AppUser, error) {
	var oAppUser domain.AppUser

	bFound, oErr := oSelf.GetById(oSelf.Index, strconv.FormatUint(uint64(iId), 10), &oAppUser)
	if oErr != nil {
		return nil, oErr
	}

	if !bFound {
		return nil, nil
	}

	return &oAppUser, nil
}

func (oSelf *AppUserModel) AddOne(oValue *domain.AppUserValue) (bool, error) {
	iId, oErr := oSelf.NextId("app_user")
	if oErr != nil {
		return false, oErr
	}

	oDoc := &domain.AppUser{Id: iId}

	if oValue.Name != nil {
		oDoc.Name = *oValue.Name
	}
	if oValue.Password != nil {
		oDoc.Password = *oValue.Password
	}
	if oValue.Balance != nil {
		oDoc.Balance = *oValue.Balance
	}

	if oErr := oSelf.IndexOne(oSelf.Index, strconv.FormatUint(uint64(iId), 10), oDoc); oErr != nil {
		return false, oErr
	}

	return true, nil
}
