package register

import (
	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"encoding/json"
	"fmt"
	"net/http"
	"time"

	bootstrap "example/bootstrap"
	container "example/container"
	utility "example/internal/utility"
	pkg "example/pkg"
	types "example/types"
)

// websocketKeys 是 Header.K 用 RSA 私鑰解開後的內容，跟 http／facade 版本
// DecryptionMiddleware／DecryptionInterceptor 解出來的 {key, iv} 是同一套格式。
type Keys struct {
	Key string `json:"key"`
	Iv  string `json:"iv"`
}

// WebsocketOnConnectFunc / WebsocketOnMessageFunc / WebsocketOnCloseFunc 是連線生命週期
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
type WebsocketOnCloseFunc func(oConn *websocket.Conn, iType int)

// WebsocketOnUnsubscribeFunc 不是 client 主動送的 application event，是 websocket
// 協定層級的斷線（iType == -1 時）：跟 OnClose 綁在同一個時間點一起觸發，讓
// 呼叫端在連線真的斷掉那一刻，順便清掉這個連線訂閱的 channel。因為觸發時 read 已經
// 失敗，沒有解析出有效的 request，簽名跟 OnClose 一樣不帶 oReq。
type WebsocketOnUnsubscribeFunc func(oConn *websocket.Conn, iType int)

// WebsocketOnHeartbeatFunc / WebsocketOnRpcFunc / WebsocketOnSubscribeFunc /
// WebsocketOnBroadcastFunc / WebsocketOnNotifyFunc / WebsocketOnRefreshFunc 都是
// 自己定義的 application event：跟 OnConnect 一樣由 client 主動送對應 event
// （heartbeat/rpc/subscribe/broadcast/notify/refresh）觸發，簽名比照 OnConnect，
// 不像 OnAuthenticate 需要回傳值決定連線生死，處理完就 continue，不會落到下面的
// onMessage。
type WebsocketOnHeartbeatFunc func(oConn *websocket.Conn, iType int, oReq *types.WebsocketRequest)
type WebsocketOnRpcFunc func(oConn *websocket.Conn, iType int, oReq *types.WebsocketRequest)
type WebsocketOnSubscribeFunc func(oConn *websocket.Conn, iType int, oReq *types.WebsocketRequest)
type WebsocketOnBroadcastFunc func(oConn *websocket.Conn, iType int, oReq *types.WebsocketRequest)
type WebsocketOnNotifyFunc func(oConn *websocket.Conn, iType int, oReq *types.WebsocketRequest)
type WebsocketOnRefreshFunc func(oConn *websocket.Conn, iType int, oReq *types.WebsocketRequest)

// WebsocketEventer 職責跟 TcpRouter 一樣：只負責「連線生命週期」機制本身
// （upgrade、read loop、斷線偵測、ping/pong keepalive），不管收到訊息／連線／斷線後
// 實際要做什麼——通訊邏輯（這支檔案）跟業務邏輯（呼叫端注入的三個 callback）完全分開。
type WebsocketEventer struct {
	upgrader       websocket.Upgrader
	onOpen         WebsocketOnOpenFunc
	onConnect      WebsocketOnConnectFunc
	onAuthenticate WebsocketOnAuthenticateFunc
	onMessage      WebsocketOnMessageFunc
	onClose        WebsocketOnCloseFunc
	onUnsubscribe  WebsocketOnUnsubscribeFunc
	onHeartbeat    WebsocketOnHeartbeatFunc
	onRpc          WebsocketOnRpcFunc
	onSubscribe    WebsocketOnSubscribeFunc
	onBroadcast    WebsocketOnBroadcastFunc
	onNotify       WebsocketOnNotifyFunc
	onRefresh      WebsocketOnRefreshFunc
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

func (oSelf *WebsocketEventer) OnClose(fnHandler WebsocketOnCloseFunc) *WebsocketEventer {
	oSelf.onClose = fnHandler
	return oSelf
}

func (oSelf *WebsocketEventer) OnUnsubscribe(fnHandler WebsocketOnUnsubscribeFunc) *WebsocketEventer {
	oSelf.onUnsubscribe = fnHandler
	return oSelf
}

func (oSelf *WebsocketEventer) OnHeartbeat(fnHandler WebsocketOnHeartbeatFunc) *WebsocketEventer {
	oSelf.onHeartbeat = fnHandler
	return oSelf
}

func (oSelf *WebsocketEventer) OnRpc(fnHandler WebsocketOnRpcFunc) *WebsocketEventer {
	oSelf.onRpc = fnHandler
	return oSelf
}

func (oSelf *WebsocketEventer) OnSubscribe(fnHandler WebsocketOnSubscribeFunc) *WebsocketEventer {
	oSelf.onSubscribe = fnHandler
	return oSelf
}

func (oSelf *WebsocketEventer) OnBroadcast(fnHandler WebsocketOnBroadcastFunc) *WebsocketEventer {
	oSelf.onBroadcast = fnHandler
	return oSelf
}

func (oSelf *WebsocketEventer) OnNotify(fnHandler WebsocketOnNotifyFunc) *WebsocketEventer {
	oSelf.onNotify = fnHandler
	return oSelf
}

func (oSelf *WebsocketEventer) OnMessage(fnHandler WebsocketOnMessageFunc) *WebsocketEventer {
	oSelf.onMessage = fnHandler
	return oSelf
}

func (oSelf *WebsocketEventer) OnRefresh(fnHandler WebsocketOnRefreshFunc) *WebsocketEventer {
	oSelf.onRefresh = fnHandler
	return oSelf
}

// ServeHTTP 讓 WebsocketEventer 可以直接掛進 http.ServeMux，用法跟其他 http.Handler 一樣。
func (oSelf *WebsocketEventer) ServeHTTP(oWriter http.ResponseWriter, oRequest *http.Request) {
	oConn, oErr := oSelf.upgrader.Upgrade(oWriter, oRequest, nil)

	if oErr != nil {
		pkg.Logger(pkg.WebsocketAdmin).Error("upgrade error", zap.Error(oErr))
		return
	}

	defer oConn.Close()

	// PONG_WAIT：連線允許完全靜默多久，超過就直接斷線。收到 Pong 會把它續回滿額——
	// 定時送出的 Ping（見下面 PING_INTERVAL）就是持續給 client 機會回 Pong，讓一直
	// 有在回應的連線不會被 PONG_WAIT 誤判斷線。
	oConn.SetReadDeadline(time.Now().Add(time.Duration(bootstrap.CONFIG.SERVICES.WEBSOCKET.PONG_WAIT) * time.Second))
	oConn.SetPongHandler(func(string) error {
		return oConn.SetReadDeadline(time.Now().Add(time.Duration(bootstrap.CONFIG.SERVICES.WEBSOCKET.PONG_WAIT) * time.Second))
	})

	if oSelf.onOpen != nil {
		oSelf.onOpen(oConn, 0)
	}

	// 定時送 Ping，跟下面讀訊息的迴圈各自獨立跑：這裡只負責照 PING_INTERVAL
	// 送 Ping，ServeHTTP 返回時 defer close(chDone) 會讓它一起停止。
	chDone := make(chan struct{})
	defer close(chDone)

	go func() {
		oTicker := time.NewTicker(time.Duration(bootstrap.CONFIG.SERVICES.WEBSOCKET.PING_INTERVAL) * time.Second)
		defer oTicker.Stop()

		for {
			select {
			case <-oTicker.C:

				oConn.WriteMessage(websocket.PingMessage, nil)
			case <-chDone:
				return
			}
		}
	}()

	// 讀取 client 訊息；handler 本身已經是 net/http 每個請求各自的 goroutine，
	// 不需要再包一層 go func()，不然這裡會直接返回，defer oConn.Close() 馬上執行，
	// 把還在等訊息的連線關掉。
	for {
		iType, aMsg, oErr := oConn.ReadMessage()

		var oWsReq types.WebsocketRequest

		if oErr == nil {

			if jsonErr := json.Unmarshal(aMsg, &oWsReq); jsonErr != nil {
				pkg.Logger(pkg.WebsocketAdmin).Error("json unmarshal error", zap.Error(jsonErr))
			}

			if (iType == 1 || iType == 2) && oSelf.onConnect != nil && oWsReq.Event == "connect" {
				oSelf.onConnect(oConn, iType, &oWsReq)
				continue
			}

			if (iType == 1 || iType == 2) && oSelf.onAuthenticate != nil && oWsReq.Event == "authenticate" {
				if !oSelf.onAuthenticate(oConn, iType, &oWsReq) {
					return
				}

				continue

			}

			if (iType == 1 || iType == 2) && oSelf.onHeartbeat != nil && oWsReq.Event == "heartbeat" {
				oSelf.onHeartbeat(oConn, iType, &oWsReq)

				continue

			}

			if (iType == 1 || iType == 2) && oSelf.onRpc != nil && oWsReq.Event == "rpc" {
				oSelf.onRpc(oConn, iType, &oWsReq)

				continue

			}

			if (iType == 1 || iType == 2) && oSelf.onSubscribe != nil && oWsReq.Event == "subscribe" {
				oSelf.onSubscribe(oConn, iType, &oWsReq)

				continue

			}

			if (iType == 1 || iType == 2) && oSelf.onBroadcast != nil && oWsReq.Event == "broadcast" {
				oSelf.onBroadcast(oConn, iType, &oWsReq)

				continue

			}

			if (iType == 1 || iType == 2) && oSelf.onNotify != nil && oWsReq.Event == "notify" {
				oSelf.onNotify(oConn, iType, &oWsReq)

				continue

			}

			if (iType == 1 || iType == 2) && oSelf.onRefresh != nil && oWsReq.Event == "refresh" {
				oSelf.onRefresh(oConn, iType, &oWsReq)

				continue

			}

			if (iType == 1 || iType == 2) && oSelf.onMessage != nil && oWsReq.Event == "message" {
				oSelf.onMessage(oConn, iType, aMsg)

				continue

			}

		}

		if oErr != nil {

			if iType == -1 {

				if oSelf.onUnsubscribe != nil {
					oSelf.onUnsubscribe(oConn, iType)
				}

				if oSelf.onClose != nil {
					oSelf.onClose(oConn, iType)
				}

				return

			}

		}

	}
}

/*

          open
		  close                                                                  完成
		  ping        / pong                                                         完成

event:

	      connect    / conntected                                  reply ✅           完成
		  heartbeat  / heartbeated                                 reply ✅           完成

		  authenticate/ authenticated                              reply ✅

--------------------------------------需要檢查是否 authenticated -----------------------------------------------


	      presence, presence-stats, history,                       reply ✅           可取消
		  rpc        / rpced                                       reply ✅

		  subscribe  / subscribed                                  reply ✅

		  chat       / chated
		  broadcast  / broadcasted                                 reply ✅ + broadcast ✅
		  message    / messaged 不需要                              reply ❌
		  notify    / notified                                     reply ✅

		  refresh     /refreshed                                  reply ✅
		  sub-refresh / sub-refreshed                             reply ✅

--------------------------------------需要檢查是否 authenticated -----------------------------------------------

		  unsubscribe / unsubscribed 不需要                        reply ❌

    method:
	value:


*/

/*
1. WebSocket 需要實現3層的 包： 協議層，應用層，動作層
(基本上就是透過 傳過來的 包 做 拆包 後 3類型路由分發)

2. 需要考慮多台 websocket server， 利用 redis PUB/SUB


│
├── A. Protocol Layer
│      └── WebSocket 協議本身定義的控制事件
│          open
│          ping
│          pong
│          close
│
├── B. Application Layer
│      └── 業務通訊語義
│          connect
│          heartbeat
│          authenticate
│          message            1 對 server, server 不回應
|          chat               1 對多 群聊
|          broadcast          1 對多 通知
│          notify             1 對 1 通知
│          rpc
│
└── C. Action / Route Layer
       └── 具體業務操作
	       broadcast:method
             ├── gift
             └── like

           notify:method
             ├── add-friend
             └── poke

           rpc :method
             ├── resource/AppUser/ShowOne
             ├── resource/AppUser/ShowOnes
             └── resource/AppUser/AddOne
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
		pkg.Logger(pkg.WebsocketAdmin).Info("OnOpen", zap.Stringer("remoteAddr", oConn.RemoteAddr()))

		sUuid := uuid.New().String()
		sPointer := fmt.Sprintf("%p", oConn)

		if _, bGotten := uuidToConnection.Get(sUuid); bGotten {
			pkg.Logger(pkg.WebsocketAdmin).Error(
				"duplicate id, disconnect",
				zap.String("uuid", sUuid),
				zap.Stringer("remoteAddr", oConn.RemoteAddr()),
			)
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
				pkg.Logger(pkg.WebsocketAdmin).Error("rsa decrypt error", zap.Error(oErr))
				return
			}

			oKeys, oErr := utility.JsonDecode[Keys](sKeys)
			if oErr != nil {
				pkg.Logger(pkg.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
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
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)

	})

	oAdminEventer.OnHeartbeat(func(oConn *websocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
		// 不用 request-id， 採用完全異步策略
		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
		}{
			Event: "heartbeated",
		})

		if oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)
	})

	oAdminEventer.OnAuthenticate(func(oConn *websocket.Conn, iType int, oWsReq *types.WebsocketRequest) bool {
		var oValue struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
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
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
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

	oAdminEventer.OnClose(func(oConn *websocket.Conn, iType int) {
		sPointer := fmt.Sprintf("%p", oConn)

		sUuid, _ := pointerToUuid.Get(sPointer)
		pointerToUuid.Del(sPointer)
		uuidToConnection.Del(sUuid)
		uuidToAuthentication.Del(sPointer)
		uuidToKeys.Del(sPointer)

		pkg.Logger(pkg.WebsocketAdmin).Info(
			"disconnected",
			zap.String("uuid", sUuid),
			zap.Stringer("remoteAddr", oConn.RemoteAddr()),
		)
	})
	oAdminEventer.OnMessage(func(oConn *websocket.Conn, iType int, aMsg []byte) {
		sPointer := fmt.Sprintf("%p", oConn)

		if bAuthenticated, _ := uuidToAuthentication.Get(sPointer); !bAuthenticated {
			pkg.Logger(pkg.WebsocketAdmin).Info("not authenticated, ignore message", zap.Stringer("remoteAddr", oConn.RemoteAddr()))
			return
		}

		oConn.WriteMessage(iType, aMsg)
	})

	oMux := http.NewServeMux()
	oMux.Handle("/Admin", oAdminEventer)

	return oMux
}
