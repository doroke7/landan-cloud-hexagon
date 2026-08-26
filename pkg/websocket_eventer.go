package pkg

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	types "example/types"
)

// WebsocketOnConnectFunc / WebsocketOnMessageFunc / WebsocketOnDisconnectFunc 是連線生命週期
// 三個時機點各自的處理方法簽名，職責跟 TcpRouter 的 method 對照表一樣：eventer 只負責在對的
// 時機呼叫對的方法，實際要做什麼交給呼叫端注入。
type WebsocketOnConnectFunc func(oConn *websocket.Conn, iType int, oReq *types.WebsocketRequest)

// WebsocketOnOpenFunc 的 iType 在 upgrade 剛完成、還沒讀過任何一個 frame 時呼叫，
// 沒有真正的 frame type 可以帶，ServeHTTP 固定傳 0，純粹是為了跟其他四個 callback 簽名一致。
type WebsocketOnOpenFunc func(oConn *websocket.Conn, iType int)

// WebsocketOnAuthenticateFunc 回傳 bool 表示驗證是否通過：true 讓連線繼續往下讀之後的訊息，
// false 讓 ServeHTTP 關閉連線——跟 onConnect／onOpen 不同，這裡的結果會影響連線生死。
type WebsocketOnAuthenticateFunc func(oConn *websocket.Conn, iType int, oReq *types.WebsocketRequest) bool
type WebsocketOnMessageFunc func(oConn *websocket.Conn, iType int, aMsg []byte)
type WebsocketOnDisconnectFunc func(oConn *websocket.Conn, iType int)

// WebsocketEventer 職責跟 TcpRouter 一樣：只負責「連線生命週期」機制本身
// （upgrade、read loop、斷線偵測），不管收到訊息／連線／斷線後實際要做什麼——
// 通訊邏輯（這支檔案）跟業務邏輯（呼叫端注入的三個 callback）完全分開。
type WebsocketEventer struct {
	upgrader       websocket.Upgrader
	onOpen         WebsocketOnOpenFunc
	onConnect      WebsocketOnConnectFunc
	onAuthenticate WebsocketOnAuthenticateFunc
	onMessage      WebsocketOnMessageFunc
	onDisconnect   WebsocketOnDisconnectFunc
}

func NewWebsocketEventer(oUpgrader websocket.Upgrader) *WebsocketEventer {
	return &WebsocketEventer{upgrader: oUpgrader}
}

func (oSelf *WebsocketEventer) OnOpen(fnHandler WebsocketOnOpenFunc) *WebsocketEventer {
	oSelf.onOpen = fnHandler
	return oSelf
}

func (oSelf *WebsocketEventer) OnConnect(fnHandler WebsocketOnConnectFunc) *WebsocketEventer {
	oSelf.onConnect = fnHandler
	return oSelf
}

func (oSelf *WebsocketEventer) OnAuthenticate(fnHandler WebsocketOnAuthenticateFunc) *WebsocketEventer {
	oSelf.onAuthenticate = fnHandler
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

	if oSelf.onOpen != nil {
		oSelf.onOpen(oConn, 0)
	}

	// 讀取 client 訊息；handler 本身已經是 net/http 每個請求各自的 goroutine，
	// 不需要再包一層 go func()，不然這裡會直接返回，defer oConn.Close() 馬上執行，
	// 把還在等訊息的連線關掉。
	for {
		iType, aMsg, oErr := oConn.ReadMessage()
		log.Println("iType:", iType)
		log.Println("aMsg:", string(aMsg))
		log.Println("oErr:", oErr)

		var oWsReq types.WebsocketRequest

		if oErr == nil {

			if jsonErr := json.Unmarshal(aMsg, &oWsReq); jsonErr != nil {
				log.Println("json unmarshal error:", jsonErr)
			}

			if (iType == 1 || iType == 2) && oSelf.onConnect != nil && oWsReq.Event == "connect" {
				oSelf.onConnect(oConn, iType, &oWsReq)
				return

			}

			if (iType == 1 || iType == 2) && oSelf.onAuthenticate != nil && oWsReq.Event == "authenticate" {
				if !oSelf.onAuthenticate(oConn, iType, &oWsReq) {
					return
				}

				continue

			}

		}

		if oErr != nil {

			if iType == -1 && oSelf.onDisconnect != nil {
				oSelf.onDisconnect(oConn, iType)
				return

			}

		}

		if oSelf.onMessage != nil {
			oSelf.onMessage(oConn, iType, aMsg)
		}
	}
}
