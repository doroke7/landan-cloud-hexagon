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
			// if oErr := oClient.Subscribe("room-01"); oErr != nil {
			// 	log.Println("subscribe error:", oErr)
			// }

			/*
			   1. 所有的 client 過來的消息 ， OnPublish 都會handler
			*/
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

			// 心跳：對應設計文件最上面 "Channel: (空)" 那組 heartbeat/heartbeat-ack。
			// RPC 命令本來就沒有 channel 概念（protocol.RPCRequest 只有 Method/Data，
			// 沒有 Channel 欄位），跟 Publish 需要帶 channel 不一樣，所以 channel="" 這種
			// 心跳/room join-leave/event 類型的協議走 RPC，不走 Publish。
			// client 端呼叫 centrifuge.rpc("heartbeat", ...)，這裡收到後直接回 heartbeat-ack，
			// 之後要記錄「最後上線時間」之類的副作用可以加在 case "heartbeat" 這裡。
			oClient.OnRPC(func(oEvent centrifuge.RPCEvent, fnCallback centrifuge.RPCCallback) {
				switch oEvent.Method {
				case "heartbeat":
					aData, _ := json.Marshal(map[string]string{"type": "rpc-ack", "method": "heartbeat"})
					fnCallback(centrifuge.RPCReply{Data: aData}, nil)
				default:
					fnCallback(centrifuge.RPCReply{}, centrifuge.ErrorMethodNotFound)
				}
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
