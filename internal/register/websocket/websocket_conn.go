package registerWebsocket

import (
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	pkg "example/pkg"
)

// outboundMessage 是塞進 Conn 寫入 channel 裡的一筆待寫資料，iType 對應
// websocket.TextMessage/BinaryMessage/PingMessage 這些 frame type。
type outboundMessage struct {
	iType int
	aData []byte
}

type WebsocketConn struct {
	Conn      *websocket.Conn
	Outbox    chan outboundMessage
	Done      chan struct{}
	CloseOnce sync.Once
}

const connOutboxSize = 32

func NewConn(oConn *websocket.Conn) *WebsocketConn {
	oSelf := &WebsocketConn{
		Conn:   oConn,
		Outbox: make(chan outboundMessage, connOutboxSize),
		Done:   make(chan struct{}),
	}

	go oSelf.runWriter()

	return oSelf
}

func (oSelf *WebsocketConn) runWriter() {
	for {
		select {
		case oMsg := <-oSelf.Outbox:
			if oErr := oSelf.Conn.WriteMessage(oMsg.iType, oMsg.aData); oErr != nil {
				pkg.Logger(pkg.WebsocketAdmin).Error("write error", zap.Error(oErr))
				return
			}
		case <-oSelf.Done:
			return
		}
	}
}

// WriteMessage 不回傳 error：實際寫入是非同步的，呼叫當下還不知道會不會成功，
// 現有呼叫端本來也都沒在檢查回傳值。
func (oSelf *WebsocketConn) WriteMessage(iType int, aData []byte) {
	select {
	case oSelf.Outbox <- outboundMessage{iType: iType, aData: aData}:
	default:
		pkg.Logger(pkg.WebsocketAdmin).Info("outbox full, drop message", zap.Int("type", iType))
	}
}

func (oSelf *WebsocketConn) ReadMessage() (int, []byte, error) {
	return oSelf.Conn.ReadMessage()
}

func (oSelf *WebsocketConn) Close() error {
	oSelf.CloseOnce.Do(func() {
		close(oSelf.Done)
	})

	return oSelf.Conn.Close()
}

func (oSelf *WebsocketConn) RemoteAddr() net.Addr {
	return oSelf.Conn.RemoteAddr()
}

func (oSelf *WebsocketConn) SetReadDeadline(oTime time.Time) error {
	return oSelf.Conn.SetReadDeadline(oTime)
}

func (oSelf *WebsocketConn) SetPongHandler(fnHandler func(string) error) {
	oSelf.Conn.SetPongHandler(fnHandler)
}
