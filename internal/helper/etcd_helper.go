package helper

import (
	"context"
	"encoding/json"
	"errors"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdHelper struct {
	*AbstractHelper
	etcd *clientv3.Client
}

func NewEtcdHelper(oAbstractHelper *AbstractHelper, oEtcd *clientv3.Client) *EtcdHelper {
	return &EtcdHelper{
		AbstractHelper: oAbstractHelper,
		etcd:           oEtcd,
	}
}

// WriteCache 把任意 value 序列化成 JSON 寫進 etcd，key 由呼叫端決定，
// 這裡只負責通用的「怎麼寫」，不管特定 domain 的 key 格式。
func (oSelf *EtcdHelper) WriteCache(sKey string, value any) error {
	sData, oErr := json.Marshal(value)
	if oErr != nil {
		return oErr
	}

	oContext := context.Background()
	_, oErr = oSelf.etcd.Put(oContext, sKey, string(sData))

	return oErr
}

// ReadCache 從 etcd 讀出 JSON 並解到 dest（傳指標進來），
// 一樣只負責通用的「怎麼讀」，key 格式跟目標型別都由呼叫端決定。
func (oSelf *EtcdHelper) ReadCache(sKey string, dest any) error {
	oContext := context.Background()
	oResponse, oErr := oSelf.etcd.Get(oContext, sKey)
	if oErr != nil {
		return oErr
	}

	if len(oResponse.Kvs) == 0 {
		oNotFoundErr := errors.New("key not found")

		return oNotFoundErr
	}

	oErr = json.Unmarshal(oResponse.Kvs[0].Value, dest)

	return oErr
}

// EvictCache 把 key 從 etcd 刪掉，用在寫入之後讓下一次讀取重新從來源撈最新資料，
// 避免寫入端自己組的資料跟實際落地的資料不一致。
func (oSelf *EtcdHelper) EvictCache(sKey string) error {
	oContext := context.Background()
	_, oErr := oSelf.etcd.Delete(oContext, sKey)

	return oErr
}
