package main

import (
	"sync"
)

type AbstractModel struct {
	appId     int
	instances sync.Map
}

func NewAbstractModel(iAppId int) *AbstractModel {
	return &AbstractModel{
		appId: iAppId,
	}
}

/*
這段代碼的問題在于， 高併發情況中 ，可能會 執行 N次，其中 有幾次執行 同時沒有讀取到，又同時間創建了
*/

func (oSelf *AbstractModel) Init(iAppId int) *AbstractModel {
	if oInstance, bOk := oSelf.instances.Load(iAppId); bOk {
		return oInstance.(*AbstractModel)
	}

	oInstance := NewAbstractModel(iAppId)
	oSelf.instances.Store(iAppId, oInstance)
	return oInstance
}
