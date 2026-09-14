package pkgNats

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/centrifugal/centrifuge"
	"github.com/nats-io/nats.go"
)

// NATSBroker 是一個最小可用的 centrifuge.Broker 實作，用 NATS core pub/sub 讓多個
// centrifuge node 之間可以互相轉發 Publication/Join/Leave 訊息，藉此讓
// cmd/centrifuge.go 可以水平擴展成多節點。刻意不支援 history/recovery——
// History 永遠回傳空、RemoveHistory 是 no-op；需要 history 的話應該換成
// RedisBroker 或另外實作，這裡只做即時轉發。
type NATSBroker struct {
	nc           *nats.Conn
	eventHandler centrifuge.BrokerEventHandler

	mu   sync.Mutex
	subs map[string]*nats.Subscription
}

var _ centrifuge.Broker = (*NATSBroker)(nil)

func NewNATSBroker(nc *nats.Conn) *NATSBroker {
	aSubs := make(map[string]*nats.Subscription)

	return &NATSBroker{
		nc:   nc,
		subs: aSubs,
	}
}

// natsBrokerMessage 是實際丟到 NATS subject 上的信封格式，用 Kind 區分
// Publication/Join/Leave，這樣一個 channel 只需要訂閱一個 subject。
type natsBrokerMessage struct {
	Kind string                 `json:"kind"`
	Data []byte                 `json:"data,omitempty"`
	Info *centrifuge.ClientInfo `json:"info,omitempty"`
	Tags map[string]string      `json:"tags,omitempty"`
}

func (oSelf *NATSBroker) subject(sChannel string) string {
	return "centrifuge." + sChannel
}

func (oSelf *NATSBroker) RegisterBrokerEventHandler(oHandler centrifuge.BrokerEventHandler) error {
	oSelf.eventHandler = oHandler
	return nil
}

func (oSelf *NATSBroker) Subscribe(sChannel string) error {
	oSelf.mu.Lock()
	defer oSelf.mu.Unlock()

	if _, bOk := oSelf.subs[sChannel]; bOk {
		return nil
	}

	sSubject := oSelf.subject(sChannel)
	oSub, err := oSelf.nc.Subscribe(sSubject, func(oMsg *nats.Msg) {
		var oMessage natsBrokerMessage
		if err := json.Unmarshal(oMsg.Data, &oMessage); err != nil {
			return
		}

		switch oMessage.Kind {
		case "join":
			_ = oSelf.eventHandler.HandleJoin(sChannel, oMessage.Info)
		case "leave":
			_ = oSelf.eventHandler.HandleLeave(sChannel, oMessage.Info)
		default:
			oPub := &centrifuge.Publication{
				Data: oMessage.Data,
				Info: oMessage.Info,
				Tags: oMessage.Tags,
			}
			_ = oSelf.eventHandler.HandlePublication(sChannel, oPub, centrifuge.StreamPosition{}, false, nil)
		}
	})
	if err != nil {
		return err
	}

	oSelf.subs[sChannel] = oSub
	return nil
}

func (oSelf *NATSBroker) Unsubscribe(sChannel string) error {
	oSelf.mu.Lock()
	defer oSelf.mu.Unlock()

	oSub, bOk := oSelf.subs[sChannel]
	if !bOk {
		return nil
	}
	delete(oSelf.subs, sChannel)
	oErr := oSub.Unsubscribe()

	return oErr
}

func (oSelf *NATSBroker) publish(sChannel string, oMessage natsBrokerMessage) error {
	aBytes, err := json.Marshal(oMessage)
	if err != nil {
		return err
	}
	sSubject := oSelf.subject(sChannel)
	oErr := oSelf.nc.Publish(sSubject, aBytes)

	return oErr
}

func (oSelf *NATSBroker) Publish(sChannel string, aData []byte, oOpts centrifuge.PublishOptions) (centrifuge.StreamPosition, bool, error) {
	err := oSelf.publish(sChannel, natsBrokerMessage{Kind: "pub", Data: aData, Info: oOpts.ClientInfo, Tags: oOpts.Tags})
	return centrifuge.StreamPosition{}, false, err
}

func (oSelf *NATSBroker) PublishJoin(sChannel string, oInfo *centrifuge.ClientInfo) error {
	oErr := oSelf.publish(sChannel, natsBrokerMessage{Kind: "join", Info: oInfo})

	return oErr
}

func (oSelf *NATSBroker) PublishLeave(sChannel string, oInfo *centrifuge.ClientInfo) error {
	oErr := oSelf.publish(sChannel, natsBrokerMessage{Kind: "leave", Info: oInfo})

	return oErr
}

func (oSelf *NATSBroker) History(sChannel string, oOpts centrifuge.HistoryOptions) ([]*centrifuge.Publication, centrifuge.StreamPosition, error) {
	return nil, centrifuge.StreamPosition{}, nil
}

func (oSelf *NATSBroker) RemoveHistory(sChannel string) error {
	return nil
}

func (oSelf *NATSBroker) Close(_ context.Context) error {
	oSelf.mu.Lock()
	defer oSelf.mu.Unlock()

	for sChannel, oSub := range oSelf.subs {
		_ = oSub.Unsubscribe()
		delete(oSelf.subs, sChannel)
	}
	return nil
}
