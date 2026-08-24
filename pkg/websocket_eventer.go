package pkg

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebsocketOnConnectFunc / WebsocketOnMessageFunc / WebsocketOnDisconnectFunc 是連線生命週期
// 三個時機點各自的處理方法簽名，職責跟 TcpRouter 的 method 對照表一樣：eventer 只負責在對的
// 時機呼叫對的方法，實際要做什麼交給呼叫端注入。
type WebsocketOnConnectFunc func(oConn *websocket.Conn)
type WebsocketOnMessageFunc func(oConn *websocket.Conn, iMessageType int, aMsg []byte)
type WebsocketOnDisconnectFunc func(oConn *websocket.Conn)

// WebsocketEventer 職責跟 TcpRouter 一樣：只負責「連線生命週期」機制本身
// （upgrade、read loop、斷線偵測），不管收到訊息／連線／斷線後實際要做什麼——
// 通訊邏輯（這支檔案）跟業務邏輯（呼叫端注入的三個 callback）完全分開。
type WebsocketEventer struct {
	upgrader     websocket.Upgrader
	onConnect    WebsocketOnConnectFunc
	onMessage    WebsocketOnMessageFunc
	onDisconnect WebsocketOnDisconnectFunc
}

func NewWebsocketEventer(oUpgrader websocket.Upgrader) *WebsocketEventer {
	return &WebsocketEventer{upgrader: oUpgrader}
}

// OnConnect 註冊連線建立完成時要執行的方法，對應 socket.io 的 .on("connect", cb)。
func (oSelf *WebsocketEventer) OnConnect(fnHandler WebsocketOnConnectFunc) *WebsocketEventer {
	oSelf.onConnect = fnHandler
	return oSelf
}

// OnMessage 註冊收到訊息時要執行的方法。
func (oSelf *WebsocketEventer) OnMessage(fnHandler WebsocketOnMessageFunc) *WebsocketEventer {
	oSelf.onMessage = fnHandler
	return oSelf
}

// OnDisconnect 註冊斷線時要執行的方法，對應 socket.io 的 .on("disconnect", cb)。
func (oSelf *WebsocketEventer) OnDisconnect(fnHandler WebsocketOnDisconnectFunc) *WebsocketEventer {
	oSelf.onDisconnect = fnHandler
	return oSelf
}

// ServeHTTP 讓 WebsocketEventer 可以直接掛進 http.ServeMux，用法跟其他 http.Handler 一樣。
func (oSelf *WebsocketEventer) ServeHTTP(oWriter http.ResponseWriter, oRequest *http.Request) {
	oConn, oErr := oSelf.upgrader.Upgrade(oWriter, oRequest, nil)

	if oErr != nil {
		log.Println(oErr)
		return
	}

	defer oConn.Close()

	if oSelf.onConnect != nil {
		oSelf.onConnect(oConn)
	}

	// 讀取 client 訊息；handler 本身已經是 net/http 每個請求各自的 goroutine，
	// 不需要再包一層 go func()，不然這裡會直接返回，defer oConn.Close() 馬上執行，
	// 把還在等訊息的連線關掉。
	for {
		iMessageType, aMsg, oErr := oConn.ReadMessage()

		if oErr != nil {
			// gorilla/websocket 讀取出錯時（連線被關閉、網路中斷...）會把 messageType
			// 設回 noFrame(-1)，用這個當作「真的斷線了」的判斷依據，而不是任何 oErr != nil 就算斷線
			if iMessageType == -1 && oSelf.onDisconnect != nil {
				oSelf.onDisconnect(oConn)
			}

			return
		}

		if oSelf.onMessage != nil {
			oSelf.onMessage(oConn, iMessageType, aMsg)
		}
	}
}
