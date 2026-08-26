package register

import (
	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"encoding/json"
	"fmt"
	"log"
	"net/http"

	bootstrap "example/bootstrap"
	container "example/container"
	utility "example/internal/utility"
	types "example/types"
)

// websocketKeys 是 Header.K 用 RSA 私鑰解開後的內容，跟 http／facade 版本
// DecryptionMiddleware／DecryptionInterceptor 解出來的 {key, iv} 是同一套格式。
type Keys struct {
	Key string `json:"key"`
	Iv  string `json:"iv"`
}

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

/*

          open                                                                        完成
		  ping        / pong

event:

	      connect    / conntected                                  reply ✅           完成
	      disconnect / disconnected 不需要                          reply ❌           完成
		  heartbeat  / heartbeated                                 reply ✅

		  authenticate/ authenticated                              reply ✅

--------------------------------------需要檢查是否 authenticated -----------------------------------------------


	      presence, presence-stats, history,                       reply ✅           可取消
		  rpc        / rpced                                       reply ✅

		  subscribe  / subscribed                                  reply ✅

		  message    / messaged 不需要                              reply ❌
		  broadcast  / broadcasted                                 reply ✅ + broadcast ✅

		  refresh     /refreshed                                  reply ✅
		  sub-refresh / sub-refreshed                             reply ✅

--------------------------------------需要檢查是否 authenticated -----------------------------------------------

		  unsubscribe / unsubscribed 不需要                        reply ❌

    method:
	value:


*/

func WebsocketInit(oContainer *container.WebsocketContainer) *http.ServeMux {

	oAdminEventer := NewWebsocketEventer(websocket.Upgrader{
		CheckOrigin: func(oRequest *http.Request) bool {
			return true
		},
	})

	var (
		pointerToUuid        = hashmap.New[string, string]()
		uuidToConnection     = hashmap.New[string, *websocket.Conn]()
		uuidToAuthentication = hashmap.New[string, bool]()
		uuidToKeys           = hashmap.New[string, Keys]()

		uuidToChannels       = hashmap.New[string, string]()
		channelToConnections = hashmap.New[string, string]()
	)

	_ = uuidToChannels
	_ = channelToConnections

	oAdminEventer.OnOpen(func(oConn *websocket.Conn, iType int) {
		log.Println("OnOpen:", oConn.RemoteAddr())

		sUuid := uuid.New().String()
		sPointer := fmt.Sprintf("%p", oConn)

		if _, bGotten := uuidToConnection.Get(sUuid); bGotten {
			log.Println("duplicate id, disconnect:", sUuid, oConn.RemoteAddr())
			oConn.Close()
			return
		}

		uuidToConnection.Set(sUuid, oConn)
		pointerToUuid.Set(sPointer, sUuid)

	})

	oAdminEventer.OnConnect(func(oConn *websocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sPointer := fmt.Sprintf("%p", oConn)

		sUuid, _ := pointerToUuid.Get(sPointer)

		if oWsReq.K != "" {
			sKeys, oErr := oContainer.RsaHelper.Decrypt(oWsReq.K, bootstrap.CONFIG.SERVICES.WEBSOCKET.ADMIN.PRIVATE_KEY)
			if oErr != nil {
				log.Println("rsa decrypt error:", oErr)
				return
			}

			oKeys, oErr := utility.JsonDecode[Keys](sKeys)
			if oErr != nil {
				log.Println("json unmarshal error:", oErr)
				return
			}

			uuidToKeys.Set(sUuid, oKeys)
		}

		aByteMessage, oErr := json.Marshal(struct {
			Event  string `json:"event"`
			Result string `json:"result"`
			RId    string `json:"r_id"`
			CId    string `json:"c_id"`
		}{
			Event:  "connected",
			Result: "",
			RId:    oWsReq.RId,
			CId:    sUuid,
		})

		if oErr != nil {
			log.Println("json marshal error:", oErr)
			return
		}

		oConn.WriteMessage(iType, aByteMessage)

	})

	oAdminEventer.OnAuthenticate(func(oConn *websocket.Conn, iType int, oWsReq *types.WebsocketRequest) bool {
		var oValue struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		fmt.Println("134 boWsReqOk=", oWsReq)

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			log.Println("json unmarshal error:", oErr)
		}

		bOk := oValue.Name == "admin" && oValue.Password == "123456"

		nCode := -1
		if bOk {
			nCode = 1
		}

		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
			Code  int    `json:"code"`
			RId   string `json:"r_id"`
		}{
			Event: "authenticated",
			Code:  nCode,
			RId:   oWsReq.RId,
		})

		if oErr != nil {
			log.Println("json marshal error:", oErr)
			return false
		}

		oConn.WriteMessage(iType, aByteMessage)

		if bOk {
			sPointer := fmt.Sprintf("%p", oConn)
			uuidToAuthentication.Set(sPointer, true)
		}

		return bOk
	})

	// disconnect 是收不到 uuid 的

	oAdminEventer.OnDisconnect(func(oConn *websocket.Conn, iType int) {
		sPointer := fmt.Sprintf("%p", oConn)

		sUuid, _ := pointerToUuid.Get(sPointer)
		pointerToUuid.Del(sPointer)
		uuidToConnection.Del(sUuid)
		uuidToAuthentication.Del(sPointer)
		uuidToKeys.Del(sPointer)

		log.Println("disconnected:", sUuid, oConn.RemoteAddr())
	})
	oAdminEventer.OnMessage(func(oConn *websocket.Conn, iType int, aMsg []byte) {
		sPointer := fmt.Sprintf("%p", oConn)

		if bAuthenticated, _ := uuidToAuthentication.Get(sPointer); !bAuthenticated {
			log.Println("not authenticated, ignore message:", oConn.RemoteAddr())
			return
		}

		oConn.WriteMessage(iType, aMsg)
	})

	oMux := http.NewServeMux()
	oMux.Handle("/Admin", oAdminEventer)

	return oMux
}
