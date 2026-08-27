package registerWebsocket

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	bootstrap "example/bootstrap"
	container "example/container"
	pkgUtility "example/pkg/utility"
	pkgWebsocket "example/pkg/websocket"
	types "example/types"
)

/*

          open
		  close                                                                      完成
		  ping        / pong                                                         完成

event:

	      connect    / conntected                                  reply ✅           完成
		  heartbeat  / heartbeated                                 reply ✅           完成

		  authenticate/ authenticated                              reply ✅           完成

--------------------------------------需要檢查是否 authenticated -----------------------------------------------


	      presence, presence-stats, history,                       reply ✅           可取消
		  rpc        / rpced                                       reply ✅           ⚠️ server 需要寫路由

		  subscribe  / subscribed                                  reply ✅           完成
		  present                                                  reply ✅           完成

		  message    / messaged 當向訊息-不需要ack                    reply ❌            完成
		  chat       / chated                                       reply ✅           完成
		  broadcast  / broadcasted                                 reply ✅ + broadcast ✅  ⚠️ server 需要依 method 寫路由（gift/like）
		  notify    / notified                                     reply ✅           ⚠️ server 需要依 method 寫路由（add-friend/poke）

		  refresh     /refreshed                                  reply ✅
		  sub-refresh / sub-refreshed                             reply ✅

--------------------------------------需要檢查是否 authenticated -----------------------------------------------

		  unsubscribe / unsubscribed                              reply ❌           完成

--------------------------------------server 不用實作 -----------------------------------------------
          publish

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
             └── nap

           notify:method
             ├── add-friend
             └── poke

           rpc :method
             ├── resource/AppUser/ShowOne
             ├── resource/AppUser/ShowOnes
             └── resource/AppUser/AddOne
*/

const (
	// wsCheckInterval / wsActiveTimeout 故意設得比正式環境短，demo 幾分鐘就能看到
	// 閒置連線被背景掃描強制關掉。
	wsCheckInterval = 5 * time.Minute
	wsActiveTimeout = 2 * time.Minute
)

func Init(oContainer *container.WebsocketContainer) *http.ServeMux {

	oHub := pkgWebsocket.NewHub()

	oAdminEventer := pkgWebsocket.NewEventer(websocket.Upgrader{
		CheckOrigin: func(oRequest *http.Request) bool {
			return true
		},
	})

	oAdminEventer.OnOpen(func(oConn *pkgWebsocket.Conn, iType int) {
		pkgUtility.Logger(pkgUtility.WebsocketAdmin).Info("OnOpen", zap.Stringer("remoteAddr", oConn.RemoteAddr()))

		if _, bOk := oHub.Add(oConn); !bOk {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error(
				"duplicate id, disconnect",
				zap.Stringer("remoteAddr", oConn.RemoteAddr()),
			)
			oConn.Close()
		}
	})

	oAdminEventer.OnPong(func(oConn *pkgWebsocket.Conn) {
		sConnectionId, _ := oHub.ConnectionId(oConn)
		oHub.Activate(sConnectionId)
	})

	oAdminEventer.OnConnect(func(oConn *pkgWebsocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sConnectionId, _ := oHub.ConnectionId(oConn)

		if oWsReq.K != "" {
			sKeys, oErr := oContainer.RsaHelper.Decrypt(oWsReq.K, bootstrap.CONFIG.SERVICES.WEBSOCKET.ADMIN.PRIVATE_KEY)
			if oErr != nil {
				pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("rsa decrypt error", zap.Error(oErr))
				return
			}

			oKeys, oErr := pkgUtility.JsonDecode[struct {
				Key string `json:"key"`
				Iv  string `json:"iv"`
			}](sKeys)
			if oErr != nil {
				pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
				return
			}

			oHub.SetKeyIv(sConnectionId, oKeys.Key, oKeys.Iv)
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
			CId:    sConnectionId,
		})

		if oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)
	})

	oAdminEventer.OnHeartbeat(func(oConn *pkgWebsocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sConnectionId, _ := oHub.ConnectionId(oConn)
		oHub.Activate(sConnectionId)

		// 不用 request-id， 採用完全異步策略
		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
		}{
			Event: "heartbeated",
		})

		if oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)
	})

	oAdminEventer.OnAuthenticate(func(oConn *pkgWebsocket.Conn, iType int, oWsReq *types.WebsocketRequest) bool {
		sConnectionId, _ := oHub.ConnectionId(oConn)

		var oValue struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
		}

		bOk := oValue.Name == "admin" && oValue.Password == "123456"

		nCode := -1
		if bOk {
			nCode = 1
		}

		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
			Code  int    `json:"code"`
			CId   string `json:"c_id"`
			RId   string `json:"r_id"`
		}{
			Event: "authenticated",
			Code:  nCode,
			CId:   sConnectionId,
			RId:   oWsReq.RId,
		})

		if oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return false
		}

		oConn.WriteMessage(iType, aByteMessage)

		if bOk {
			// NOTE: 目前沒有真的 admin user id，先用 cId 當識別綁上去。
			oHub.Authenticate(sConnectionId, sConnectionId)
		}

		return bOk
	})

	// OnRpc 先把 Method 取出來，但目前還沒有真的依它分派到對應的業務邏輯
	// （resource/AppUser/ShowOne 之類的 Action Route），只回 rpced 確認收到。
	oAdminEventer.OnRpc(func(oConn *pkgWebsocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sMethod := oWsReq.Method
		pkgUtility.Logger(pkgUtility.WebsocketAdmin).Info("OnRpc", zap.String("method", sMethod))

		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
			RId   string `json:"r_id"`
		}{
			Event: "rpced",
			RId:   oWsReq.RId,
		})

		if oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)
	})

	oAdminEventer.OnSubscribe(func(oConn *pkgWebsocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sConnectionId, _ := oHub.ConnectionId(oConn)

		if !oHub.Authenticated(sConnectionId) {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Info("not authenticated, ignore subscribe", zap.String("cid", sConnectionId))
			return
		}

		var oValue struct {
			Channel string `json:"channel"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
			return
		}

		oHub.JoinConnectionIdChannel(sConnectionId, oValue.Channel)

		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
			CId   string `json:"c_id"`
			RId   string `json:"r_id"`
			Value struct {
				Channel string `json:"channel"`
			} `json:"value"`
		}{
			Event: "subscribed",
			CId:   sConnectionId,
			RId:   oWsReq.RId,
			Value: struct {
				Channel string `json:"channel"`
			}{
				Channel: oValue.Channel,
			},
		})

		if oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)

		// 廣播給這個 channel 目前所有訂閱的 cid（含剛加入的這條連線自己），讓大家
		// 知道有新成員加入；跟 OnChat 一樣直接重複用同一份 aByteMessage，不用
		// 另外組一份不同的 payload。
		go oHub.PublishToChannel(oValue.Channel, iType, aByteMessage)
	})

	oAdminEventer.OnPresent(func(oConn *pkgWebsocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sConnectionId, _ := oHub.ConnectionId(oConn)

		var oValue struct {
			Channel string `json:"channel"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
			return
		}

		aAdminUserIds := oHub.AdminUsersInChannel(oValue.Channel)
		aOnes := make([]struct {
			Id string `json:"id"`
		}, 0, len(aAdminUserIds))

		for _, sAdminUserId := range aAdminUserIds {
			aOnes = append(aOnes, struct {
				Id string `json:"id"`
			}{Id: sAdminUserId})
		}

		aByteMessage, oErr := json.Marshal(struct {
			Event  string `json:"event"`
			CId    string `json:"c_id"`
			RId    string `json:"r_id"`
			Result struct {
				Channel string `json:"channel"`
				Ones    []struct {
					Id string `json:"id"`
				} `json:"ones"`
			} `json:"result"`
		}{
			Event: "presented",
			CId:   sConnectionId,
			RId:   oWsReq.RId,
			Result: struct {
				Channel string `json:"channel"`
				Ones    []struct {
					Id string `json:"id"`
				} `json:"ones"`
			}{
				Channel: oValue.Channel,
				Ones:    aOnes,
			},
		})

		if oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)
	})

	oAdminEventer.OnChat(func(oConn *pkgWebsocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sConnectionId, _ := oHub.ConnectionId(oConn)

		var oValue struct {
			Channel string `json:"channel"`
			Text    string `json:"text"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
			return
		}

		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
			CId   string `json:"c_id"`
			RId   string `json:"r_id"`
		}{
			Event: "chated",
			CId:   sConnectionId,
			RId:   oWsReq.RId,
		})

		if oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)

		if len(oHub.ConnectionIdsInChannel(oValue.Channel)) == 0 {
			return
		}

		aByteChat, oErr := json.Marshal(struct {
			Event  string `json:"event"`
			CId    string `json:"c_id"`
			RId    string `json:"r_id"`
			Result struct {
				Channel string `json:"channel"`
				Text    string `json:"text"`
			} `json:"result"`
		}{
			Event: "chat",
			CId:   sConnectionId,
			RId:   oWsReq.RId,
			Result: struct {
				Channel string `json:"channel"`
				Text    string `json:"text"`
			}{
				Channel: oValue.Channel,
				Text:    oValue.Text,
			},
		})

		if oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		go oHub.PublishToChannel(oValue.Channel, iType, aByteChat)
	})

	// OnBroadcast 目前還沒有真的依 Method（gift/like 之類）分派到對應的業務
	// 邏輯，先把整個 value 原封不動連同 method 一起轉發給頻道成員，讓前端自己
	// 依 method 處理內容；跟 OnChat 一樣，channel 沒人訂閱就直接忽略。
	oAdminEventer.OnBroadcast(func(oConn *pkgWebsocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sConnectionId, _ := oHub.ConnectionId(oConn)

		var oValue struct {
			Channel string `json:"channel"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
			return
		}

		aByteMessage, oErr := json.Marshal(struct {
			Event  string `json:"event"`
			Method string `json:"method"`
			CId    string `json:"c_id"`
			RId    string `json:"r_id"`
			Value  struct {
				Channel string `json:"channel"`
			} `json:"value"`
		}{
			Event:  "broadcasted",
			Method: oWsReq.Method,
			CId:    sConnectionId,
			RId:    oWsReq.RId,
			Value: struct {
				Channel string `json:"channel"`
			}{
				Channel: oValue.Channel,
			},
		})

		if oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)

		if len(oHub.ConnectionIdsInChannel(oValue.Channel)) == 0 {
			return
		}

		aByteBroadcast, oErr := json.Marshal(struct {
			Event  string          `json:"event"`
			Method string          `json:"method"`
			CId    string          `json:"c_id"`
			RId    string          `json:"r_id"`
			Value  json.RawMessage `json:"value"`
		}{
			Event:  "broadcast",
			Method: oWsReq.Method,
			CId:    sConnectionId,
			RId:    oWsReq.RId,
			Value:  oWsReq.Value,
		})

		if oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		go oHub.PublishToChannel(oValue.Channel, iType, aByteBroadcast)
	})

	oAdminEventer.OnNotify(func(oConn *pkgWebsocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sConnectionId, _ := oHub.ConnectionId(oConn)

		var oValue struct {
			AdminUserId string `json:"admin_user_id"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
			return
		}

		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
			CId   string `json:"c_id"`
			RId   string `json:"r_id"`
		}{
			Event: "notified",
			CId:   sConnectionId,
			RId:   oWsReq.RId,
		})

		if oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)

		if len(oHub.AdminUserIds(oValue.AdminUserId)) == 0 {
			return
		}

		aByteNotify, oErr := json.Marshal(struct {
			Event  string          `json:"event"`
			CId    string          `json:"c_id"`
			RId    string          `json:"r_id"`
			Method string          `json:"method"`
			Value  json.RawMessage `json:"value"`
		}{
			Event:  "notify",
			CId:    sConnectionId,
			RId:    oWsReq.RId,
			Method: oWsReq.Method,
			Value:  oWsReq.Value,
		})

		if oErr != nil {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		go oHub.PublishToAdminUserId(oValue.AdminUserId, iType, aByteNotify)
	})

	oAdminEventer.OnUnsubscribe(func(oConn *pkgWebsocket.Conn, iType int) {
		sConnectionId, _ := oHub.ConnectionId(oConn)
		oHub.LeaveConnectionId(sConnectionId)
	})

	oAdminEventer.OnMessage(func(oConn *pkgWebsocket.Conn, iType int, aMsg []byte) {
		sConnectionId, _ := oHub.ConnectionId(oConn)

		if !oHub.Authenticated(sConnectionId) {
			pkgUtility.Logger(pkgUtility.WebsocketAdmin).Info("not authenticated, ignore message", zap.Stringer("remoteAddr", oConn.RemoteAddr()))
			return
		}

		// DO NOTHING ，不回傳 ack
	})

	oAdminEventer.OnClose(func(oConn *pkgWebsocket.Conn, iType int) {
		sConnectionId := oHub.Remove(oConn)

		pkgUtility.Logger(pkgUtility.WebsocketAdmin).Info(
			"disconnected",
			zap.String("cid", sConnectionId),
			zap.Stringer("remoteAddr", oConn.RemoteAddr()),
		)
	})

	go oHub.SweepCron(wsCheckInterval, wsActiveTimeout, nil)

	oMux := http.NewServeMux()
	oMux.Handle("/Admin", oAdminEventer)

	return oMux
}
