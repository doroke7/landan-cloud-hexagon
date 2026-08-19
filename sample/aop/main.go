package aop

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"example/bootstrap"
)

/**
  Go 1.26.4 版本 不支持 struct method 在method 上賦予泛型 --> 只能用 局部函數
  Go 1.26.4 版本 不支持 把函數當作 泛型 -->

*/

/**

1. 核心邏輯對照表 (Core Logic Mapping)
+----------------+---------------+-----------------------+----------------+
| 註解 (Annotation)| 核心動作      | 對接邏輯 (Data Flow)  | 適用場景 (Get) |
+----------------+---------------+-----------------------+----------------+
| @Cacheable     | 讀取 / 存入   | 有就拿，沒有就存。    | GetUser (查詢) |
+----------------+---------------+-----------------------+----------------+
| @CachePut      | 更新 / 覆蓋   | 不管怎樣都執行並更新。| UpdateUser(修改)|
+----------------+---------------+-----------------------+----------------+
| @CacheEvict    | 刪除 / 清空   | 直接把緩存踢出去。    | DeleteUser(刪除)|
+----------------+---------------+-----------------------+----------------+


*/

/*


sample/aop 的用法：


func GetUser(id int64) (User, error) {
      oUser, oErr := pkg.Cacheable("", 10*time.Minute, func() (User, error) {
              oUser, oErr := queryUserFromDB(id) // 真正查 DB
              if oErr != nil {
                      return oUser, oErr
              }

              return oUser, nil
      }, id) // id 會被拼進 cache key（Md5(aParams...)）
      if oErr != nil {
              return oUser, oErr
      }

      return oUser, nil
}



pkg/aop 的用法：

func (oSelf *UserService) GetUser(oCtx context.Context, nID int64) (User, error) {
      var oUser User
      sKey := fmt.Sprintf("user:%d", nID)

      oErr := oSelf.oAop.Cacheable(oCtx, sKey, 10*time.Minute, &oUser, func(oCtx context.Context) (interface{}, error) {
              oUser, oErr := queryUserFromDB(nID) // 真正查 DB
              if oErr != nil {
                      return oUser, oErr
              }

              return oUser, nil
      })
      if oErr != nil {
              return oUser, oErr
      }

      return oUser, nil
}


兩個方式的差異是：
1. pkg package 函數 vs struct 函數
2. 不使用 context vs 使用 context
3. 需要指定型別 （User）在注入函數 vs 直接注入結構體（ &oUser）

*. 結論 下面的寫法比較優雅
*/

type Aop struct {
	Redis   redis.UniversalClient
	Context context.Context
}

func NewAop() *Aop {
	oRedis, oErr := bootstrap.NewRedis()
	if oErr != nil {
		log.Fatalf("aop: failed to init redis: %v", oErr)
	}

	return &Aop{
		Redis:   oRedis,
		Context: context.Background(),
	}
}

var DefaultApp = sync.OnceValue(NewAop)

func callerKey(nSkip int) string {
	oPC, _, _, bOk := runtime.Caller(nSkip)
	if !bOk {
		return "unknown"
	}
	oFn := runtime.FuncForPC(oPC)
	if oFn == nil {
		return "unknown"
	}
	sName := oFn.Name()
	if nIdx := strings.LastIndex(sName, "/"); nIdx >= 0 {
		sName = sName[nIdx+1:]
	}
	if nIdx := strings.Index(sName, "."); nIdx >= 0 {
		sName = sName[nIdx+1:]
	}
	sName = strings.NewReplacer("(*", "", ")", "", ".", ":").Replace(sName)
	return sName
}

func Md5(aParams ...any) string {
	aData, _ := json.Marshal(aParams)
	aHash := md5.Sum(aData)
	return fmt.Sprintf("%x", aHash)
}

// Cacheable 類似 Hyperf #[Cacheable]
// cFn 為無參數閉包，aParams 用於產生 MD5 cache key
func Cacheable[G any](sKey string, oTtl time.Duration, cFn func() (G, error), aParams ...any) (G, error) {
	oAop := DefaultApp()
	var oModel G
	sCacheKey := sKey
	if sCacheKey == "" {
		sCacheKey = callerKey(2)
	}
	if len(aParams) > 0 {
		sCacheKey += ":" + Md5(aParams...)
	}

	aRaw, oErr := oAop.Redis.Get(oAop.Context, sCacheKey).Bytes()
	if oErr == nil {
		var oResult G
		if oJsonErr := json.Unmarshal(aRaw, &oResult); oJsonErr == nil {
			return oResult, nil
		}
	}

	oResult, oErr := cFn()
	if oErr != nil {
		return oModel, oErr
	}
	if aData, oJsonErr := json.Marshal(oResult); oJsonErr == nil {
		oAop.Redis.Set(oAop.Context, sCacheKey, aData, oTtl)
	}
	return oResult, nil
}

// CachePut 類似 Hyperf #[CachePut]
func CachePut[G any](sKey string, oTtl time.Duration, cFn func() (G, error), aParams ...any) (G, error) {
	oAop := DefaultApp()
	var oModel G
	sCacheKey := sKey
	if sCacheKey == "" {
		sCacheKey = callerKey(2)
	}
	if len(aParams) > 0 {
		sCacheKey += ":" + Md5(aParams...)
	}

	oResult, oErr := cFn()
	if oErr != nil {
		return oModel, oErr
	}
	if aData, oJsonErr := json.Marshal(oResult); oJsonErr == nil {
		oAop.Redis.Set(oAop.Context, sCacheKey, aData, oTtl)
	}
	return oResult, nil
}

// CacheEvict 類似 Hyperf #[CacheEvict]
func CacheEvict(sKey string, cFn func() error, aParams ...any) error {
	oAop := DefaultApp()
	sCacheKey := sKey
	if sCacheKey == "" {
		sCacheKey = callerKey(2)
	}
	if len(aParams) > 0 {
		sCacheKey += ":" + Md5(aParams...)
	}

	if oErr := cFn(); oErr != nil {
		return oErr
	}
	oAop.Redis.Del(oAop.Context, sCacheKey)
	return nil
}

// CacheEvictByPattern 用 pattern 批次刪除 cache（如 "user:*"）
func CacheEvictByPattern(sPattern string, cFn func() error) error {
	oAop := DefaultApp()
	if oErr := cFn(); oErr != nil {
		return oErr
	}

	var nCursor uint64
	for {
		aKeys, nNext, oErr := oAop.Redis.Scan(oAop.Context, nCursor, sPattern, 100).Result()
		if oErr != nil {
			break
		}
		if len(aKeys) > 0 {
			oAop.Redis.Del(oAop.Context, aKeys...)
		}
		nCursor = nNext
		if nCursor == 0 {
			break
		}
	}
	return nil
}
