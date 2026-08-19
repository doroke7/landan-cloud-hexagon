package pkg

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	bootstrap "example/bootstrap"
)

// Aop 結構體：不帶泛型，是一個純粹、萬用的切面工具箱
type Aop struct {
	Cache   redis.UniversalClient
	Context context.Context
}

// NewAop 初始化時注入萬用的 Redis 客戶端，以及程序等級的全局 ctx。
// 全局 ctx 取代每次呼叫方傳進來的 ctx：程序收到中斷/終止訊號時，
// 所有 Cacheable／CacheEvict 都會跟著一起停止，不受個別呼叫方有沒有接好 ctx 影響。
func NewAop(oContext context.Context, oCache redis.UniversalClient) *Aop {
	return &Aop{
		Cache:   oCache,
		Context: oContext,
	}
}

// Cacheable 順利綁定為 Aop 結構體的方法！
// 💡 關鍵設計：
// 1. pDest: 傳入用來接收結果的「結構體指針」（如同 json.Unmarshal 的第二個參數）
// 2. cFn: 核心業務函數（如查 DB），回傳萬用的 interface{}
func (oSelf *Aop) Cacheable(sKey string, oTtl time.Duration, pDest interface{}, cFn func() (interface{}, error)) error {

	oCurrentContext, cCancel := context.WithTimeout(
		oSelf.Context,
		time.Duration(bootstrap.CONFIG.REDIS.TIMEOUT)*time.Millisecond,
	)
	defer cCancel()

	// 1. 前置安全檢查
	if oErr := oCurrentContext.Err(); oErr != nil {
		return oErr
	}

	// 2. 前置切面 (Before)：查快取
	sJsonStr, oErr := oSelf.Cache.Get(oCurrentContext, sKey).Result()
	if oErr == nil {
		// 🎉 快取命中 (Hit)！直接反序列化進傳進來的 pDest 指針
		if oErr := json.Unmarshal([]byte(sJsonStr), pDest); oErr == nil {
			fmt.Printf("🎯 [Aop] 完美命中快取！資料已自動注入。Key: '%s'\n", sKey)
			return nil
		}
	} else if !errors.Is(oErr, redis.Nil) {
		// 快取降級：Redis 異常不中斷主業務，讓流量穿透去查 DB
		fmt.Printf("⚠️ [Aop] 快取連線異常: %v，自動降級穿透。\n", oErr)
	}

	fmt.Printf("🔍 [Aop] 快取未命中 (Miss)，準備執行核心業務... Key: '%s'\n", sKey)

	// 3. 執行核心業務 (Join Point：例如查 MySQL)
	oRes, oErr := cFn()
	if oErr != nil {
		return oErr
	}

	// 4. 後置切面 (After Returning)：將真實業務回傳的物件，序列化同步到 Redis
	aByteData, oErr := json.Marshal(oRes)
	if oErr == nil {
		oSelf.Cache.Set(oCurrentContext, sKey, string(aByteData), oTtl)
		fmt.Printf("💾 [Aop] 數據已自動同步至 Redis（TTL: %v）。Key: '%s'\n", oTtl, sKey)
	}

	// 5. 💡 核心魔法：如果快取沒中，查完 DB 後，要把結果深拷貝（Deep Copy）給外面的 pDest 指針
	if oErr := json.Unmarshal(aByteData, pDest); oErr != nil {
		return oErr
	}

	return nil
}

// CacheEvict 先執行核心業務 (Join Point：例如更新/刪除 DB)，成功後把對應的快取踢掉
func (oSelf *Aop) CacheEvict(sKey string, cFn func() error) error {

	oCurrentContext, cCancel := context.WithTimeout(
		oSelf.Context,
		time.Duration(bootstrap.CONFIG.REDIS.TIMEOUT)*time.Millisecond,
	)
	defer cCancel()

	// 1. 前置安全檢查
	if oErr := oCurrentContext.Err(); oErr != nil {
		return oErr
	}

	// 2. 執行核心業務 (Join Point：例如更新 MySQL)
	if oErr := cFn(); oErr != nil {
		return oErr
	}

	// 3. 後置切面 (After Returning)：把快取踢出去
	if oErr := oSelf.Cache.Del(oCurrentContext, sKey).Err(); oErr != nil && !errors.Is(oErr, redis.Nil) {
		fmt.Printf("⚠️ [Aop] 快取清除失敗: %v，Key: '%s'\n", oErr, sKey)
		return oErr
	}

	fmt.Printf("🗑️ [Aop] 快取已清除。Key: '%s'\n", sKey)
	return nil
}

// Key 把任意個數的參數轉成字串、用 ":" 合併，再取 md5，
// 方便呼叫端組出固定長度、不用自己擔心特殊字元的 cache key。
func (oSelf *Aop) Key(sPrefix string, aArgs ...interface{}) string {
	aParts := make([]string, len(aArgs))
	for i, oArg := range aArgs {
		aParts[i] = fmt.Sprint(oArg)
	}

	aSum := md5.Sum([]byte(strings.Join(aParts, ":")))
	return sPrefix + ":" + hex.EncodeToString(aSum[:])
}

// Ttl 在 oBase 上下加減一個 [0, oBase/2] 的隨機值，
// 用來讓大量 key 的到期時間錯開，避免同時失效造成 cache stampede。
func (oSelf *Aop) Ttl(oBase time.Duration) time.Duration {
	oHalf := oBase / 2
	iJitter := rand.Int64N(int64(oHalf) + 1)

	if rand.IntN(2) == 0 {
		return oBase + time.Duration(iJitter)
	}
	return oBase - time.Duration(iJitter)
}
