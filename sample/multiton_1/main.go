package main

import (
	"fmt"
	"sync"
)

// AbstractModel 你的基礎結構
type AbstractModel struct {
	AppID int
	// ... 其他業務欄位
}

// 承載器：這是多例模式的核心「坑位」
type instanceHolder struct {
	instance *AbstractModel
	once     sync.Once // 每個實例專屬的初始化鎖
}

// ModelManager 管理所有的多例
type ModelRepository struct {
	instanceHolders sync.Map // Key: int (AppID), Value: *instanceHolder
}

func (oSelf *ModelRepository) InitModel(iAppId int) *AbstractModel {
	// 1. 原子佔位：嘗試獲取或建立一個承載器
	// LoadOrStore 保證了針對同一個 iAppId，大家拿到的 holder 指針是同一個
	oValue, _ := oSelf.instanceHolders.LoadOrStore(iAppId, &instanceHolder{})
	oInstanceHoler := oValue.(*instanceHolder)

	// 2. 執行初始化：利用承載器內部的 once 確保 New 只執行一次
	oInstanceHoler.once.Do(func() {
		fmt.Printf("[系統] 正在為 AppID:%d 建立唯一的物理實例...\n", iAppId)
		// 這裡執行昂貴的初始化，例如連資料庫或讀取設定
		oInstanceHoler.instance = &AbstractModel{
			AppID: iAppId,
		}
	})

	return oInstanceHoler.instance
}
