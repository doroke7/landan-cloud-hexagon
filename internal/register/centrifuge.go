package register

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/centrifugal/centrifuge"

	container "example/container"
)

/*
   Channel:
   Type:  rpc /rpc-ack
   Method: heartbeat
   Value:

   Channel:
   Type: rpc / rpc-ack
   Method: Admin/Authentication/Authenticator/SignIn
   Value: {}

   Channel:
   Type: room / room-ack
   Method: join/leave
   Value: room-01

   Channel: room-01
   Type: message / message-ack
   Method: broast
   Value: {}

   Channel: room-01
   Type: message / messages-ack
   Method: send
   Value: {}

*/

// CentrifugeInit 組裝 centrifuge.Node 跟它的 websocket handler，回傳 oNode 是因為
// cmd 那邊還需要用它做 Shutdown、以及背景 ticker 呼叫 oNode.Publish(...) 廣播。
func CentrifugeInit(oContainer *container.CentrifugeContainer) (*centrifuge.Node, http.Handler) {

	oNode, oErr := centrifuge.New(
		centrifuge.Config{
			LogLevel: centrifuge.LogLevelDebug,
		},
	)
	if oErr != nil {
		panic(oErr)
	}

	oNode.SetBroker(oContainer.NatsBroker)

	// Connecting：驗證/接受連線，這個 demo 不做任何驗證，一律接受成匿名連線。
	// Credentials 一定要給值（UserID 留空即代表匿名），不然 centrifuge 會在
	// connectCmd 判斷 credentials == nil 直接以 bad request 斷線。
	oNode.OnConnecting(
		func(oCtx context.Context, oEvent centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
			return centrifuge.ConnectReply{
				Credentials: &centrifuge.Credentials{UserID: ""},
			}, nil
		},
	)

	// Connect：連線建立完成後才能拿到 *centrifuge.Client，在這裡掛 client 層級的事件、訂閱 channel。
	oNode.OnConnect(
		func(oClient *centrifuge.Client) {

			// 斷線
			oClient.OnDisconnect(func(oEvent centrifuge.DisconnectEvent) {})

			if oErr := oClient.Subscribe("all"); oErr != nil {
				log.Println("subscribe error:", oErr)
			}

			oClient.OnPublish(func(oEvent centrifuge.PublishEvent, fnCallback centrifuge.PublishCallback) {
				var oPayload struct {
					Method string          `json:"method"`
					Value  json.RawMessage `json:"value"`
				}
				if oErr := json.Unmarshal(oEvent.Data, &oPayload); oErr != nil {
					log.Println("publish payload 解析失敗:", oErr)
				} else {
					log.Printf(
						"收到訊息 channel=%s method=%s value=%s\n",
						oEvent.Channel,
						oPayload.Method,
						oPayload.Value,
					)
				}

				fnCallback(centrifuge.PublishReply{}, nil)
			})

			oClient.OnRPC(func(oEvent centrifuge.RPCEvent, fnCallback centrifuge.RPCCallback) {
				switch oEvent.Method {
				case "heartbeat":
					aData, _ := json.Marshal(map[string]string{"type": "rpc-ack", "method": "heartbeat"})
					fnCallback(centrifuge.RPCReply{Data: aData}, nil)
				default:
					fnCallback(centrifuge.RPCReply{}, centrifuge.ErrorMethodNotFound)
				}
			})

			// OnRefresh：client-side 連線過期時要不要延長連線。
			oClient.OnRefresh(func(oEvent centrifuge.RefreshEvent, fnCallback centrifuge.RefreshCallback) {
				fnCallback(centrifuge.RefreshReply{}, centrifuge.ErrorNotAvailable)
			})

			// OnMessage：client 用 Send（單向訊息，沒有 reply）送過來的資料，沒有 callback 可以回。
			oClient.OnMessage(func(oEvent centrifuge.MessageEvent) {
				//
			})

			// OnSubRefresh：client-side 訂閱過期時要不要延長訂閱。
			oClient.OnSubRefresh(func(oEvent centrifuge.SubRefreshEvent, fnCallback centrifuge.SubRefreshCallback) {
				fnCallback(centrifuge.SubRefreshReply{}, centrifuge.ErrorNotAvailable)
			})

			// OnSubscribe：client 端主動 subscribe 任意 channel 的權限判斷，"all"/"room-01"
			// 目前都是上面 server-side 主動 Subscribe，不會走到這個 handler，先占位不開放。
			oClient.OnSubscribe(func(oEvent centrifuge.SubscribeEvent, fnCallback centrifuge.SubscribeCallback) {
				fnCallback(centrifuge.SubscribeReply{}, centrifuge.ErrorNotAvailable)
			})

			// OnUnsubscribe：純通知事件，沒有 callback 可以回。
			oClient.OnUnsubscribe(func(oEvent centrifuge.UnsubscribeEvent) {
				//
			})

			// OnPresence / OnPresenceStats：channel 在線名單／人數查詢。
			oClient.OnPresence(func(oEvent centrifuge.PresenceEvent, fnCallback centrifuge.PresenceCallback) {
				fnCallback(centrifuge.PresenceReply{}, centrifuge.ErrorNotAvailable)
			})
			oClient.OnPresenceStats(func(oEvent centrifuge.PresenceStatsEvent, fnCallback centrifuge.PresenceStatsCallback) {
				fnCallback(centrifuge.PresenceStatsReply{}, centrifuge.ErrorNotAvailable)
			})

			// OnHistory：channel 歷史訊息查詢。
			oClient.OnHistory(func(oEvent centrifuge.HistoryEvent, fnCallback centrifuge.HistoryCallback) {
				fnCallback(centrifuge.HistoryReply{}, centrifuge.ErrorNotAvailable)
			})
		},
	)

	// 啟動 node
	if oErr := oNode.Run(); oErr != nil {
		panic(oErr)
	}

	fnHandler := centrifuge.NewWebsocketHandler(
		oNode,
		centrifuge.WebsocketConfig{
			// CORS：centrifuge 預設用 sameHostOriginCheck，要求 Origin host 跟 request Host
			// 一致，跨源的瀏覽器前端會直接被拒絕，這裡跟 pkg/websocket_router.go 同一套慣例全部放行。
			CheckOrigin: func(oRequest *http.Request) bool { return true },
		},
	)

	oMux := http.NewServeMux()
	oMux.Handle("/connection/websocket", fnHandler)

	return oNode, oMux
}
