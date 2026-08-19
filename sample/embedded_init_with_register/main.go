package main

import (
	"fmt"
	"sync"
)

type IModel interface {
}

type AbstractModel struct {
	instances sync.Map
	factories sync.Map // 存放註冊的工廠 func(int) IModel
}

func (oSelf *AbstractModel) Register(name string, fFactory func(int) IModel) {
	oSelf.factories.Store(name, fFactory)
}

type instanceHolder struct {
	instance IModel
	once     sync.Once
}

// Init 接受一個傳入的 New 函數
func (oSelf *AbstractModel) Init(name string, iAppId int) IModel {
	cacheKey := fmt.Sprintf("%s_%d", name, iAppId)

	// 1. 原子佔位：獲取或創建該 ID 的 Holder
	val, _ := oSelf.instances.LoadOrStore(cacheKey, &instanceHolder{})
	holder := val.(*instanceHolder)

	// 2. 原子施工：利用 Once 確保 factory 只跑一次
	holder.once.Do(func() {
		if f, ok := oSelf.factories.Load(name); ok {
			fmt.Printf("--- [系統] 第一次建立 %s (ID: %d) ---\n", name, iAppId)
			factory := f.(func(int) IModel)
			holder.instance = factory(iAppId)
		}
	})

	return holder.instance
}

///////////////////////////////////////////////////////////////////

type AppUserModel struct {
	*AbstractModel
	AppId    int
	Nickname string
}

// 這就是你說的 NewAppUserModel
func NewAppUserModel(id int) IModel {
	return &AppUserModel{
		AppId:    id,
		Nickname: "Landan_Makati",
	}
}

type OrderModel struct {
	*AbstractModel
	OrderId int
}

// 這就是 NewOrderModel
func NewOrderModel(id int) IModel {
	return &OrderModel{OrderId: id}
}

var oAbstractModel = &AbstractModel{}

func init() {

	oAbstractModel.Register("AppUser", NewAppUserModel)
	oAbstractModel.Register("OrderUser", NewOrderModel)

}

func main() {
	// 假設有一個全域管理員
	// 1. 獲取用戶：傳入 ID 和 構造函數
	oUserModel := oAbstractModel.Init("AppUser", 1).(*AppUserModel)
	oOrderModel := oAbstractModel.Init("Order", 1).(*OrderModel)

	fmt.Printf("用戶: %v", oUserModel)
	fmt.Printf("訂單: %v", oOrderModel)

}
