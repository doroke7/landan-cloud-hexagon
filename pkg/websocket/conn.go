package pkgWebsocket

import (
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	pkgUtility "example/pkg/utility"
)

// outboundMessage 是塞進 Conn 寫入 channel 裡的一筆待寫資料，iType 對應
// websocket.TextMessage/BinaryMessage/PingMessage 這些 frame type。
type outboundMessage struct {
	iType int
	aData []byte
}

type Conn struct {
	Conn      *websocket.Conn
	Outbox    chan outboundMessage
	Done      chan struct{}
	CloseOnce sync.Once
}

const connOutboxSize = 32

func NewConn(oConn *websocket.Conn) *Conn {
	oSelf := &Conn{
		Conn:   oConn,
		Outbox: make(chan outboundMessage, connOutboxSize),
		Done:   make(chan struct{}),
	}

	go oSelf.runWriter()

	return oSelf
}

func (oSelf *Conn) runWriter() {
	for {
		select {
		case oMsg := <-oSelf.Outbox:
			if oErr := oSelf.Conn.WriteMessage(oMsg.iType, oMsg.aData); oErr != nil {
				pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("write error", zap.Error(oErr))
				return
			}
		case <-oSelf.Done:
			return
		}
	}
}

// WriteMessage 不回傳 error：實際寫入是非同步的，呼叫當下還不知道會不會成功，
// 現有呼叫端本來也都沒在檢查回傳值。
func (oSelf *Conn) WriteMessage(iType int, aData []byte) {
	select {
	case oSelf.Outbox <- outboundMessage{iType: iType, aData: aData}:
	default:
		pkgUtility.Logger(pkgUtility.WebsocketAdmin).Info("outbox full, drop message", zap.Int("type", iType))
	}
}

func (oSelf *Conn) ReadMessage() (int, []byte, error) {
	return oSelf.Conn.ReadMessage()
}

func (oSelf *Conn) Close() error {
	oSelf.CloseOnce.Do(func() {
		close(oSelf.Done)
	})

	return oSelf.Conn.Close()
}

func (oSelf *Conn) RemoteAddr() net.Addr {
	return oSelf.Conn.RemoteAddr()
}

func (oSelf *Conn) SetReadDeadline(oTime time.Time) error {
	return oSelf.Conn.SetReadDeadline(oTime)
}

func (oSelf *Conn) SetPongHandler(fnHandler func(string) error) {
	oSelf.Conn.SetPongHandler(fnHandler)
}
