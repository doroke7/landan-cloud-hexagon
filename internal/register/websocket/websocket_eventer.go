package registerWebsocket

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	bootstrap "example/bootstrap"
	pkgUtility "example/pkg/utility"
	types "example/types"
)

// WebsocketOnConnectFunc / WebsocketOnMessageFunc / WebsocketOnCloseFunc 是連線生命週期
// 三個時機點各自的處理方法簽名，職責跟 TcpRouter 的 method 對照表一樣：eventer 只負責在對的
// 時機呼叫對的方法，實際要做什麼交給呼叫端注入。
type WebsocketOnConnectFunc func(oConn *WebsocketConn, iType int, oReq *types.WebsocketRequest)

// WebsocketOnOpenFunc 的 iType 在 upgrade 剛完成、還沒讀過任何一個 frame 時呼叫，
// 沒有真正的 frame type 可以帶，ServeHTTP 固定傳 0，純粹是為了跟其他四個 callback 簽名一致。
type WebsocketOnOpenFunc func(oConn *WebsocketConn, iType int)

// WebsocketOnAuthenticateFunc 回傳 bool 表示驗證是否通過：true 讓連線繼續往下讀之後的訊息，
// false 讓 ServeHTTP 關閉連線——跟 onConnect／onOpen 不同，這裡的結果會影響連線生死。
type WebsocketOnAuthenticateFunc func(oConn *WebsocketConn, iType int, oReq *types.WebsocketRequest) bool
type WebsocketOnMessageFunc func(oConn *WebsocketConn, iType int, aMsg []byte)
type WebsocketOnCloseFunc func(oConn *WebsocketConn, iType int)

// WebsocketOnUnsubscribeFunc 有兩種觸發時機：一是 client 主動送 event: "unsubscribe"；
// 二是 websocket 協定層級的斷線（iType == -1 時），這種情況跟 OnClose 綁在同一個時間點
// 一起觸發，讓呼叫端在連線真的斷掉那一刻，順便清掉這個連線訂閱的 channel。兩種情境都
// 不需要 reply、也不需要新的 request 資料（清的是已經記錄住的訂閱狀態），所以共用同一個
// 簽名，不帶 oReq。
type WebsocketOnUnsubscribeFunc func(oConn *WebsocketConn, iType int)

// WebsocketOnHeartbeatFunc / WebsocketOnRpcFunc / WebsocketOnSubscribeFunc /
// WebsocketOnBroadcastFunc / WebsocketOnNotifyFunc / WebsocketOnRefreshFunc /
// WebsocketOnChatFunc / WebsocketOnPresentFunc 都是自己定義的 application
// event：跟 OnConnect 一樣由 client 主動送對應 event（heartbeat/rpc/subscribe/
// broadcast/notify/refresh/chat/present）觸發，簽名比照 OnConnect，不像
// OnAuthenticate 需要回傳值決定連線生死，處理完就 continue，不會落到下面的
// onMessage。
type WebsocketOnHeartbeatFunc func(oConn *WebsocketConn, iType int, oReq *types.WebsocketRequest)
type WebsocketOnRpcFunc func(oConn *WebsocketConn, iType int, oReq *types.WebsocketRequest)
type WebsocketOnSubscribeFunc func(oConn *WebsocketConn, iType int, oReq *types.WebsocketRequest)
type WebsocketOnBroadcastFunc func(oConn *WebsocketConn, iType int, oReq *types.WebsocketRequest)
type WebsocketOnNotifyFunc func(oConn *WebsocketConn, iType int, oReq *types.WebsocketRequest)
type WebsocketOnRefreshFunc func(oConn *WebsocketConn, iType int, oReq *types.WebsocketRequest)
type WebsocketOnChatFunc func(oConn *WebsocketConn, iType int, oReq *types.WebsocketRequest)
type WebsocketOnPresentFunc func(oConn *WebsocketConn, iType int, oReq *types.WebsocketRequest)

// WebsocketOnPongFunc 是 websocket 協定層級的 Pong（client 回應 ServeHTTP 定時送出的
// Ping）：觸發時機在 gorilla 的 SetPongHandler 裡，沒有 frame type、也沒有解析出
// request，簽名只帶 oConn，讓呼叫端可以用它更新連線的最後活躍時間。
type WebsocketOnPongFunc func(oConn *WebsocketConn)

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
	onPong         WebsocketOnPongFunc
	onChat         WebsocketOnChatFunc
	onPresent      WebsocketOnPresentFunc
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

func (oSelf *WebsocketEventer) OnPong(fnHandler WebsocketOnPongFunc) *WebsocketEventer {
	oSelf.onPong = fnHandler
	return oSelf
}

func (oSelf *WebsocketEventer) OnChat(fnHandler WebsocketOnChatFunc) *WebsocketEventer {
	oSelf.onChat = fnHandler
	return oSelf
}

func (oSelf *WebsocketEventer) OnPresent(fnHandler WebsocketOnPresentFunc) *WebsocketEventer {
	oSelf.onPresent = fnHandler
	return oSelf
}

// ServeHTTP 讓 WebsocketEventer 可以直接掛進 http.ServeMux，用法跟其他 http.Handler 一樣。
func (oSelf *WebsocketEventer) ServeHTTP(oWriter http.ResponseWriter, oRequest *http.Request) {
	oRawConn, oErr := oSelf.upgrader.Upgrade(oWriter, oRequest, nil)

	if oErr != nil {
		pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("upgrade error", zap.Error(oErr))
		return
	}

	oConn := NewConn(oRawConn)

	defer oConn.Close()

	// PONG_WAIT：連線允許完全靜默多久，超過就直接斷線。收到 Pong 會把它續回滿額——
	// 定時送出的 Ping（見下面 PING_INTERVAL）就是持續給 client 機會回 Pong，讓一直
	// 有在回應的連線不會被 PONG_WAIT 誤判斷線。
	oConn.SetReadDeadline(time.Now().Add(time.Duration(bootstrap.CONFIG.SERVICES.WEBSOCKET.PONG_WAIT) * time.Second))
	oConn.SetPongHandler(func(string) error {
		if oSelf.onPong != nil {
			oSelf.onPong(oConn)
		}

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
				pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json unmarshal error", zap.Error(jsonErr))
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

			if (iType == 1 || iType == 2) && oSelf.onChat != nil && oWsReq.Event == "chat" {
				oSelf.onChat(oConn, iType, &oWsReq)

				continue

			}

			if (iType == 1 || iType == 2) && oSelf.onPresent != nil && oWsReq.Event == "present" {
				oSelf.onPresent(oConn, iType, &oWsReq)

				continue

			}

			if (iType == 1 || iType == 2) && oSelf.onMessage != nil && oWsReq.Event == "message" {
				oSelf.onMessage(oConn, iType, aMsg)

				continue

			}

			// client 主動送 event: "unsubscribe"，跟下面 iType == -1（協定層級斷線）
			// 共用同一個 onUnsubscribe：兩者都不需要 oReq（清掉的是已經記錄在
			// uuidToChannels 裡的訂閱狀態，不是這次訊息帶來的新資料），reply 也都
			// 不需要，所以簽名維持 (oConn, iType) 就夠兩種情境共用。
			if (iType == 1 || iType == 2) && oSelf.onUnsubscribe != nil && oWsReq.Event == "unsubscribe" {
				oSelf.onUnsubscribe(oConn, iType)

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
