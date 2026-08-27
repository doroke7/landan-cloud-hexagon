package registerCentrifuge

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/centrifugal/centrifuge"
	"go.uber.org/zap"

	container "example/container"
	pkgUtility "example/pkg/utility"
)

// subRefreshTTL 故意設很短，是為了讓 demo 能在幾秒內就看到 client-side subscription
// refresh 真的觸發，不用等真的過期時間（通常是幾十分鐘/幾小時）那麼久。
const subRefreshTTL = 8 * time.Second

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

// Init 組裝 centrifuge.Node 跟它的 websocket handler，回傳 oNode 是因為
// cmd 那邊還需要用它做 Shutdown、以及背景 ticker 呼叫 oNode.Publish(...) 廣播。
func Init(oContainer *container.CentrifugeContainer) (*centrifuge.Node, http.Handler) {

	oNode, oErr := centrifuge.New(
		centrifuge.Config{
			LogLevel: centrifuge.LogLevelDebug,
		},
	)
	if oErr != nil {
		panic(oErr)
	}

	oNode.SetBroker(oContainer.NatsBroker)

	oNode.OnConnecting(func(oCtx context.Context, oEvent centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		return centrifuge.ConnectReply{
			Credentials: &centrifuge.Credentials{UserID: ""},
		}, nil
	},
	)

	// Connect：連線建立完成後才能拿到 *centrifuge.Client，在這裡掛 client 層級的事件、訂閱 channel。
	oNode.OnConnect(func(oClient *centrifuge.Client) {

		// 斷線
		oClient.OnDisconnect(func(oEvent centrifuge.DisconnectEvent) {
			pkgUtility.Logger(pkgUtility.Centrifuge).Info("client OnDisconnect")

		})

		if oErr := oClient.Subscribe("all"); oErr != nil {
			pkgUtility.Logger(pkgUtility.Centrifuge).Error("subscribe error", zap.Error(oErr))
		}

		// OnPresence ~= 特殊 RPC,  client -> server -> client 一次雙向 返回
		// 有 ack
		oClient.OnPresence(func(oEvent centrifuge.PresenceEvent, fnCallback centrifuge.PresenceCallback) {
			fnCallback(centrifuge.PresenceReply{
				Result: &centrifuge.PresenceResult{
					Presence: map[string]*centrifuge.ClientInfo{
						"fake-client-01": {ClientID: "fake-client-01", UserID: "user-01"},
						"fake-client-02": {ClientID: "fake-client-02", UserID: "user-02"},
						"fake-client-03": {ClientID: "fake-client-03", UserID: ""},
					},
				},
			}, nil)
		})

		// OnPresenceStats ~= 特殊 RPC,  client -> server -> client 一次雙向 返回
		// 有 ack
		oClient.OnPresenceStats(func(oEvent centrifuge.PresenceStatsEvent, fnCallback centrifuge.PresenceStatsCallback) {
			fnCallback(centrifuge.PresenceStatsReply{
				Result: &centrifuge.PresenceStatsResult{
					PresenceStats: centrifuge.PresenceStats{NumClients: 101, NumUsers: 101},
				},
			}, nil)
		})

		// OnHistory ~= 特殊 RPC,  client -> server -> client 一次雙向 返回
		// 有 ack
		oClient.OnHistory(func(oEvent centrifuge.HistoryEvent, fnCallback centrifuge.HistoryCallback) {
			aPublications := []*centrifuge.Publication{
				{Offset: 1, Data: []byte(`{"method":"chat_message","value":"假歷史訊息 #1"}`)},
				{Offset: 2, Data: []byte(`{"method":"chat_message","value":"假歷史訊息 #2"}`)},
				{Offset: 3, Data: []byte(`{"method":"chat_message","value":"假歷史訊息 #3"}`)},
				{Offset: 4, Data: []byte(`{"method":"chat_message","value":"假歷史訊息 #4"}`)},
			}

			fnCallback(centrifuge.HistoryReply{
				Result: &centrifuge.HistoryResult{
					StreamPosition: centrifuge.StreamPosition{Offset: 4, Epoch: "fake"},
					Publications:   aPublications,
				},
			}, nil)
		})

		// OnMessage： client -> server 單向 不返回
		// 無 ack
		oClient.OnMessage(func(oEvent centrifuge.MessageEvent) {
			pkgUtility.Logger(pkgUtility.Centrifuge).Info("client OnMessage", zap.ByteString("data", oEvent.Data))
		})

		//                            -> client
		// OnPublish client -> server -> broadcast
		// 同時有 ack
		oClient.OnPublish(func(oEvent centrifuge.PublishEvent, fnCallback centrifuge.PublishCallback) {
			pkgUtility.Logger(pkgUtility.Centrifuge).Info("client OnPublish", zap.ByteString("data", oEvent.Data))

			var oPayload struct {
				Method string          `json:"method"`
				Value  json.RawMessage `json:"value"`
			}
			if oErr := json.Unmarshal(oEvent.Data, &oPayload); oErr != nil {
				pkgUtility.Logger(pkgUtility.Centrifuge).Error("publish payload 解析失敗", zap.Error(oErr))
			} else {
				pkgUtility.Logger(pkgUtility.Centrifuge).Info(
					"收到訊息",
					zap.String("channel", oEvent.Channel),
					zap.String("method", oPayload.Method),
					zap.ByteString("value", oPayload.Value),
				)
			}

			fnCallback(centrifuge.PublishReply{}, nil)
		})

		// OnSubscribe client.js Subscribe 後呼叫
		oClient.OnSubscribe(func(oEvent centrifuge.SubscribeEvent, fnCallback centrifuge.SubscribeCallback) {
			pkgUtility.Logger(pkgUtility.Centrifuge).Info("client OnSubscribe 訂閱", zap.String("channel", oEvent.Channel))
			fnCallback(centrifuge.SubscribeReply{
				Options: centrifuge.SubscribeOptions{
					ExpireAt: time.Now().Unix() + int64(subRefreshTTL.Seconds()),
				},
				ClientSideRefresh: true,
			}, nil)
		})

		// OnRPC client -> server -> client 一次雙向 返回
		// 有 ack
		oClient.OnRPC(func(oEvent centrifuge.RPCEvent, fnCallback centrifuge.RPCCallback) {
			switch oEvent.Method {
			case "heartbeat":
				aData, _ := json.Marshal(map[string]string{"type": "rpc-ack", "method": "heartbeat"})
				fnCallback(centrifuge.RPCReply{Data: aData}, nil)
			default:
				fnCallback(centrifuge.RPCReply{}, centrifuge.ErrorMethodNotFound)
			}
		})

		// OnRefresh： connect 後自動 ，套件內核會 定時自動呼叫刷新。
		oClient.OnRefresh(func(oEvent centrifuge.RefreshEvent, fnCallback centrifuge.RefreshCallback) {
			// log.Printf("client OnRefresh")

			fnCallback(centrifuge.RefreshReply{}, centrifuge.ErrorNotAvailable)
		})

		// OnSubRefresh client.js Subscribe 一個channel 後，套件內核會 定時自動呼叫刷新。
		oClient.OnSubRefresh(func(oEvent centrifuge.SubRefreshEvent, fnCallback centrifuge.SubRefreshCallback) {
			// log.Printf("client OnSubRefresh 刷新 channel=%s\n", oEvent.Channel)

			fnCallback(centrifuge.SubRefreshReply{
				ExpireAt: time.Now().Unix() + int64(subRefreshTTL.Seconds()),
			}, nil)
		})

		// OnUnsubscribe client F5關閉瀏覽器 -> 觸發發生。
		oClient.OnUnsubscribe(func(oEvent centrifuge.UnsubscribeEvent) {
			pkgUtility.Logger(pkgUtility.Centrifuge).Info(
				"client OnUnsubscribe 取消訂閱",
				zap.String("channel", oEvent.Channel),
				zap.Uint32("code", oEvent.Unsubscribe.Code),
				zap.String("reason", oEvent.Unsubscribe.Reason),
				zap.Bool("serverSide", oEvent.ServerSide),
			)
		})

	})

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
