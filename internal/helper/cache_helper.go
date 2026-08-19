package helper

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type CacheHelper struct {
	*AbstractHelper
	redis redis.UniversalClient
}

func NewCacheHelper(oAbstractHelper *AbstractHelper, oRedis redis.UniversalClient) *CacheHelper {
	return &CacheHelper{
		AbstractHelper: oAbstractHelper,
		redis:          oRedis,
	}
}

// WriteCache 把任意 value 序列化成 JSON 寫進 redis，key 由呼叫端決定，
// 這裡只負責通用的「怎麼寫」，不管特定 domain 的 key 格式。
func (oSelf *CacheHelper) WriteCache(sKey string, value any) error {
	sData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return oSelf.redis.Set(context.Background(), sKey, sData, 0).Err()
}

// ReadCache 從 redis 讀出 JSON 並解到 dest（傳指標進來），
// 一樣只負責通用的「怎麼讀」，key 格式跟目標型別都由呼叫端決定。
func (oSelf *CacheHelper) ReadCache(sKey string, dest any) error {
	sData, err := oSelf.redis.Get(context.Background(), sKey).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(sData), dest)
}

// EvictCache 把 key 從 redis 刪掉，用在寫入之後讓下一次讀取重新從來源撈最新資料，
// 避免寫入端自己組的資料跟實際落地的資料不一致（例如 auto increment ID 沒帶回來）。
func (oSelf *CacheHelper) EvictCache(sKey string) error {
	return oSelf.redis.Del(context.Background(), sKey).Err()
}
