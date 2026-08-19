package main

import (
	"fmt"
	"sync"
)

// IModel 確保所有子類都有基本的行為能力
type IModel interface {
	GetAppId() int
}

// instanceHolder 每個實例專屬的「保險箱」
type instanceHolder struct {
	instance IModel
	once     sync.Once // 核心：保證初始化只執行一次
}

// AbstractModel 父類：全局註冊中心與緩存管理器
type AbstractModel struct {
	factories sync.Map // 存儲註冊的工廠函數 key: string, value: func(int) IModel
	instances sync.Map // 存儲已生成的實例持有者 key: string, value: *instanceHolder
}

// 全域唯一的管理對象
var oModelManager = &AbstractModel{}

// Register 子類主動呼叫，完成「生產線」登記
func (oSelf *AbstractModel) Register(name string, factory func(int) IModel) {
	oSelf.factories.Store(name, factory)
}

// Init 核心邏輯：高併發安全的初始化與獲取
func (oSelf *AbstractModel) Init(name string, iAppId int) IModel {
	cacheKey := fmt.Sprintf("%s_%d", name, iAppId)

	// 1. sync.Map 原子操作：獲取或創建一個「坑位」
	// 這裡完全不需要手動 Lock/Unlock
	val, _ := oSelf.instances.LoadOrStore(cacheKey, &instanceHolder{})
	holder := val.(*instanceHolder)

	// 2. sync.Once 確保只有第一個進來的請求會執行 factory
	holder.once.Do(func() {
		if f, ok := oSelf.factories.Load(name); ok {
			fmt.Printf("[系統啟動] 正在初始化 %s (AppId: %d)...\n", name, iAppId)
			factory := f.(func(int) IModel)
			holder.instance = factory(iAppId)
		}
	})

	return holder.instance
}

// ---------------------------------------------------------
// 子類實作：AppUserModel
// ---------------------------------------------------------

type AppUserModel struct {
	AppId    int
	Nickname string
}

func (u *AppUserModel) GetAppId() int { return u.AppId }

func init() {
	// 子類主動註冊：告訴父類如果遇到 "AppUser" 該怎麼 New 出來
	oModelManager.Register("AppUser", func(id int) IModel {
		return &AppUserModel{
			AppId:    id,
			Nickname: "Landan_Backend", // 模擬數據載入
		}
	})
}

// ---------------------------------------------------------
// 執行測試
// ---------------------------------------------------------
/*

簡單的說， 這個思維，透過預先的註冊
AppUserModel
OrderModel 的產生函數，
之後需要 new 一個新實例的時候 都可以統一 的地方產生
*/
func main() {
	var wg sync.WaitGroup

	// 模擬 10 個併發請求同時獲取同一個實例
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// 獲取並斷言回子類
			user := oModelManager.Init("AppUser", 888).(*AppUserModel)
			fmt.Printf("併發請求 %d 拿到指針: %p, 名字: %s\n", id, user, user.Nickname)
		}(i)
	}

	wg.Wait()
}
