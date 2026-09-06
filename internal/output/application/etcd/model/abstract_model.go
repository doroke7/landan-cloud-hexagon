package outputApplicationEtcdModel

import (
	outputApplicationEtcd "example/internal/output/application/etcd"
)

// AbstractEtcd 由 container wire 注入，EtcdHelper / Context 提升上去；
// model 這層目前不需要額外欄位，只是把 etcd 的共用資源包一層。
type AbstractModel struct {
	*outputApplicationEtcd.AbstractEtcd
}

func NewAbstractModel(oAbstractEtcd *outputApplicationEtcd.AbstractEtcd) *AbstractModel {
	return &AbstractModel{
		AbstractEtcd: oAbstractEtcd,
	}
}
