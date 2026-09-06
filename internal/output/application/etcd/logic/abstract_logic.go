package outputApplicationEtcdLogic

import (
	outputApplicationEtcd "example/internal/output/application/etcd"
)

// AbstractEtcd 由 container wire 注入，EtcdHelper / Context 提升上去；
// logic 這層目前不需要額外欄位，只是把 etcd 的共用資源包一層。
type AbstractLogic struct {
	*outputApplicationEtcd.AbstractEtcd
}

func NewAbstractLogic(oAbstractEtcd *outputApplicationEtcd.AbstractEtcd) *AbstractLogic {
	return &AbstractLogic{
		AbstractEtcd: oAbstractEtcd,
	}
}
