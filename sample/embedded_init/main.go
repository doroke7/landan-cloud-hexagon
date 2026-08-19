package main

import (
	"fmt"
	"sync"
)

type IModel interface {
}

type AbstractModel struct {
	instances sync.Map
}

// Init 接受一個傳入的 New 函數
func (oSelf *AbstractModel) Init(iAppId int, fConstructor func(int) IModel) IModel {
	// 1. 嘗試從 Map 載入
	if oInstance, bOk := oSelf.instances.Load(iAppId); bOk {
		return oInstance.(IModel)
	}

	// 2. 沒找到，直接執行你傳進來的 NewAppUserModel 或 NewOrderModel
	oNewInstance := fConstructor(iAppId)

	// 3. 存入並回傳
	oSelf.instances.Store(iAppId, oNewInstance)
	return oNewInstance
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

func main() {
	// 假設有一個全域管理員
	oAbstractModel := &AbstractModel{}

	// 1. 獲取用戶：傳入 ID 和 構造函數
	oUserModel := oAbstractModel.Init(1, NewAppUserModel).(*AppUserModel)

	// 2. 獲取訂單：傳入 ID 和 構造函數
	oOrderModel := oAbstractModel.Init(1001, NewOrderModel).(*OrderModel)

	fmt.Printf("用戶: %s, 訂單ID: %d\n", oUserModel.Nickname, oOrderModel.OrderId)
}
