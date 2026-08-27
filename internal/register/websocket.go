package register

import (
	"sync/atomic"
	"time"

	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"encoding/json"
	"fmt"
	"net/http"

	bootstrap "example/bootstrap"
	container "example/container"
	utility "example/internal/utility"
	pkg "example/pkg"
	types "example/types"
)

// Session 是一條連線目前已知的所有狀態，取代原本四個各自獨立、卻都用同一個
// cid 當 key 的 hashmap.Map（cidToConnection／cidToAuthentication／cidToKeys／
// cidToActivedAt）。Session 本身整個是不可變的值，包含 Connection 在內——
// 不嵌套任何內層 struct，也沒有自己的鎖或原子欄位。cidToSession 存的是
// *atomic.Pointer[Session] 這個「格子」：格子的位置固定不變，格子裡指向
// 哪一份 Session 快照才會變。要更新哪個欄位，就整份複製、改掉那個欄位、
// 透過 atomic.Pointer 把格子整個換成新快照，讀的一方 Load() 拿到的永遠是
// 同一個時間點、完整一致的快照。
//
// 這條連線自己的 read loop（OnPong／OnHeartbeat／OnAuthenticate／
// OnConnect）跟另一個獨立跑的定時逾時掃描 goroutine 都會讀這個格子，
// 但只有前者會寫；各個 handler 裡都是直接「Load 舊快照、複製、改欄位、
// 整個 Store 換新」，這種寫法只在單一寫入者時才安全——如果之後有其他
// goroutine 也要寫，得改成 CompareAndSwap 迴圈才不會遺失更新。
type Session struct {
	Connection    *pkg.Conn
	ActivedAt     time.Time
	Authenticated bool
	Key           string
	Iv            string
}

func NewSession(oConn *pkg.Conn) *atomic.Pointer[Session] {
	oSession := new(atomic.Pointer[Session])
	oSession.Store(&Session{Connection: oConn, ActivedAt: time.Now()})

	return oSession
}

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

func WebsocketInit(oContainer *container.WebsocketContainer) *http.ServeMux {

	oAdminEventer := pkg.NewWebsocketEventer(websocket.Upgrader{
		CheckOrigin: func(oRequest *http.Request) bool {
			return true
		},
	})

	var (
		pointerToCid = hashmap.New[string, string]()
		cidToSession = hashmap.New[string, *atomic.Pointer[Session]]()

		auIdCIds     = pkg.NewBiMultiMap[string, string]()
		cIdsChannels = pkg.NewBiMultiMap[string, string]()
	)

	oAdminEventer.OnOpen(func(oConn *pkg.Conn, iType int) {
		pkg.Logger(pkg.WebsocketAdmin).Info("OnOpen", zap.Stringer("remoteAddr", oConn.RemoteAddr()))

		sCId := uuid.New().String()
		sPointer := fmt.Sprintf("%p", oConn)

		if _, bGotten := cidToSession.Get(sCId); bGotten {
			pkg.Logger(pkg.WebsocketAdmin).Error(
				"duplicate id, disconnect",
				zap.String("cid", sCId),
				zap.Stringer("remoteAddr", oConn.RemoteAddr()),
			)
			oConn.Close()
			return
		}

		cidToSession.Set(sCId, NewSession(oConn))
		pointerToCid.Set(sPointer, sCId)

	})

	oAdminEventer.OnPong(func(oConn *pkg.Conn) {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

		if oSession, bGotten := cidToSession.Get(sCId); bGotten {
			oNew := *oSession.Load()
			oNew.ActivedAt = time.Now()
			oSession.Store(&oNew)
		}
	})

	oAdminEventer.OnConnect(func(oConn *pkg.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sPointer := fmt.Sprintf("%p", oConn)

		sCId, _ := pointerToCid.Get(sPointer)

		if oWsReq.K != "" {
			sKeys, oErr := oContainer.RsaHelper.Decrypt(oWsReq.K, bootstrap.CONFIG.SERVICES.WEBSOCKET.ADMIN.PRIVATE_KEY)
			if oErr != nil {
				pkg.Logger(pkg.WebsocketAdmin).Error("rsa decrypt error", zap.Error(oErr))
				return
			}

			oKeys, oErr := utility.JsonDecode[struct {
				Key string `json:"key"`
				Iv  string `json:"iv"`
			}](sKeys)
			if oErr != nil {
				pkg.Logger(pkg.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
				return
			}

			if oSession, bGotten := cidToSession.Get(sCId); bGotten {
				oNew := *oSession.Load()
				oNew.Key = oKeys.Key
				oNew.Iv = oKeys.Iv
				oSession.Store(&oNew)
			}
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
			CId:    sCId,
		})

		if oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)

	})

	oAdminEventer.OnHeartbeat(func(oConn *pkg.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

		if oSession, bGotten := cidToSession.Get(sCId); bGotten {
			oNew := *oSession.Load()
			oNew.ActivedAt = time.Now()
			oSession.Store(&oNew)
		}

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

	oAdminEventer.OnAuthenticate(func(oConn *pkg.Conn, iType int, oWsReq *types.WebsocketRequest) bool {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

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
			CId   string `json:"c_id"`
			RId   string `json:"r_id"`
		}{
			Event: "authenticated",
			Code:  nCode,
			CId:   sCId,
			RId:   oWsReq.RId,
		})

		if oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return false
		}

		oConn.WriteMessage(iType, aByteMessage)

		if bOk {
			if oSession, bGotten := cidToSession.Get(sCId); bGotten {
				oNew := *oSession.Load()
				oNew.Authenticated = true
				oSession.Store(&oNew)
			}

			auIdCIds.Insert(sCId, sCId)

			cIdsChannels.Insert(sCId, "/")
		}

		return bOk
	})

	// OnRpc 先把 Method 取出來，但目前還沒有真的依它分派到對應的業務邏輯
	// （resource/AppUser/ShowOne 之類的 Action Route），只回 rpced 確認收到。
	oAdminEventer.OnRpc(func(oConn *pkg.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sMethod := oWsReq.Method
		pkg.Logger(pkg.WebsocketAdmin).Info("OnRpc", zap.String("method", sMethod))

		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
			RId   string `json:"r_id"`
		}{
			Event: "rpced",
			RId:   oWsReq.RId,
		})

		if oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)
	})

	oAdminEventer.OnSubscribe(func(oConn *pkg.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

		if len(auIdCIds.Right(sCId)) == 0 {
			pkg.Logger(pkg.WebsocketAdmin).Info("not authenticated, ignore subscribe", zap.String("cid", sCId))
			return
		}

		var oValue struct {
			Channel string `json:"channel"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
			return
		}

		cIdsChannels.Insert(sCId, oValue.Channel)
		fmt.Println("356 sCId=", sCId)
		fmt.Println("356 oValue.Channel=", oValue.Channel)

		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
			CId   string `json:"c_id"`
			RId   string `json:"r_id"`
			Value struct {
				Channel string `json:"channel"`
			} `json:"value"`
		}{
			Event: "subscribed",
			CId:   sCId,
			RId:   oWsReq.RId,
			Value: struct {
				Channel string `json:"channel"`
			}{
				Channel: oValue.Channel,
			},
		})

		if oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)

		// 廣播給這個 channel 目前所有訂閱的 cid（含剛加入的這條連線自己），讓大家
		// 知道有新成員加入；跟 OnChat 一樣直接重複用同一份 aByteMessage，不用
		// 另外組一份不同的 payload。
		go func() {
			for _, sTargetCId := range cIdsChannels.Right(oValue.Channel) {
				oTargetSession, bGotten := cidToSession.Get(sTargetCId)
				if !bGotten {
					continue
				}

				oTargetSession.Load().Connection.WriteMessage(iType, aByteMessage)
			}
		}()
	})

	oAdminEventer.OnPresent(func(oConn *pkg.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

		var oValue struct {
			Channel string `json:"channel"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
			return
		}
		aTargetCIds := cIdsChannels.Right(oValue.Channel)

		fmt.Println("411 aTargetCIds=", aTargetCIds)

		oSeenAuIds := make(map[string]struct{}, len(aTargetCIds))
		aOnes := make([]struct {
			Id string `json:"id"`
		}, 0, len(aTargetCIds))

		for _, sTargetCId := range aTargetCIds {
			for _, sTargetAuId := range auIdCIds.Right(sTargetCId) {
				if _, bSeen := oSeenAuIds[sTargetAuId]; bSeen {
					continue
				}
				oSeenAuIds[sTargetAuId] = struct{}{}

				aOnes = append(aOnes, struct {
					Id string `json:"id"`
				}{Id: sTargetAuId})
			}
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
			CId:   sCId,
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
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)
	})

	oAdminEventer.OnChat(func(oConn *pkg.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

		var oValue struct {
			Channel string `json:"channel"`
			Text    string `json:"text"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
			return
		}

		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
			CId   string `json:"c_id"`
			RId   string `json:"r_id"`
		}{
			Event: "chated",
			CId:   sCId,
			RId:   oWsReq.RId,
		})

		if oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)

		aTargetCIds := cIdsChannels.Right(oValue.Channel)
		if len(aTargetCIds) == 0 {
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
			CId:   sCId,
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
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}
		go func() {
			for _, sTargetCId := range aTargetCIds {
				oTargetSession, bGotten := cidToSession.Get(sTargetCId)
				if !bGotten {
					continue
				}

				oTargetSession.Load().Connection.WriteMessage(iType, aByteChat)
			}
		}()

	})

	// OnBroadcast 目前還沒有真的依 Method（gift/like 之類）分派到對應的業務
	// 邏輯，先把整個 value 原封不動連同 method 一起轉發給頻道成員，讓前端自己
	// 依 method 處理內容；跟 OnChat 一樣，channel 沒人訂閱就直接忽略。
	oAdminEventer.OnBroadcast(func(oConn *pkg.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

		var oValue struct {
			Channel string `json:"channel"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
			return
		}

		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
			CId   string `json:"c_id"`
			RId   string `json:"r_id"`
		}{
			Event: "broadcasted",
			CId:   sCId,
			RId:   oWsReq.RId,
		})

		if oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)

		aTargetCIds := cIdsChannels.Right(oValue.Channel)
		if len(aTargetCIds) == 0 {
			return
		}

		aByteBroadcast, oErr := json.Marshal(struct {
			Event  string          `json:"event"`
			CId    string          `json:"c_id"`
			RId    string          `json:"r_id"`
			Method string          `json:"method"`
			Value  json.RawMessage `json:"value"`
		}{
			Event:  "broadcast",
			CId:    sCId,
			RId:    oWsReq.RId,
			Method: oWsReq.Method,
			Value:  oWsReq.Value,
		})

		if oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		go func() {
			for _, sTargetCId := range aTargetCIds {
				oTargetSession, bGotten := cidToSession.Get(sTargetCId)
				if !bGotten {
					continue
				}

				oTargetSession.Load().Connection.WriteMessage(iType, aByteBroadcast)
			}
		}()
	})

	oAdminEventer.OnNotify(func(oConn *pkg.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

		var oValue struct {
			AdminUserId string `json:"admin_user_id"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
			return
		}

		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
			CId   string `json:"c_id"`
			RId   string `json:"r_id"`
		}{
			Event: "notified",
			CId:   sCId,
			RId:   oWsReq.RId,
		})

		if oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)

		aTargetCIds := auIdCIds.Left(oValue.AdminUserId)
		if len(aTargetCIds) == 0 {
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
			CId:    sCId,
			RId:    oWsReq.RId,
			Method: oWsReq.Method,
			Value:  oWsReq.Value,
		})

		if oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		go func() {
			for _, sTargetCId := range aTargetCIds {
				oTargetSession, bGotten := cidToSession.Get(sTargetCId)
				if !bGotten {
					continue
				}

				oTargetSession.Load().Connection.WriteMessage(iType, aByteNotify)
			}
		}()
	})

	oAdminEventer.OnUnsubscribe(func(oConn *pkg.Conn, iType int) {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

		fmt.Println("679 我 OnUnsubscribe了, sCId=", sCId)

		cIdsChannels.RemoveLeft(sCId)
	})

	oAdminEventer.OnMessage(func(oConn *pkg.Conn, iType int, aMsg []byte) {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

		oSession, bGotten := cidToSession.Get(sCId)
		if !bGotten || !oSession.Load().Authenticated {
			pkg.Logger(pkg.WebsocketAdmin).Info("not authenticated, ignore message", zap.Stringer("remoteAddr", oConn.RemoteAddr()))
			return
		}

		// DO NOTHING ，不回傳 ack
	})

	// OnClose 是收不到 cid 的
	oAdminEventer.OnClose(func(oConn *pkg.Conn, iType int) {
		sPointer := fmt.Sprintf("%p", oConn)

		sCId, _ := pointerToCid.Get(sPointer)
		pointerToCid.Del(sPointer)
		cidToSession.Del(sCId)
		fmt.Println("704 我 OnClose, sCId=", sCId)

		cIdsChannels.RemoveLeft(sCId)
		auIdCIds.RemoveRight(sCId)

		pkg.Logger(pkg.WebsocketAdmin).Info(
			"disconnected",
			zap.String("cid", sCId),
			zap.Stringer("remoteAddr", oConn.RemoteAddr()),
		)
	})
	go func() {
		const (
			checkInterval = 5 * time.Minute
			activeTimeout = 2 * time.Minute
		)

		oTicker := time.NewTicker(checkInterval)
		defer oTicker.Stop()

		for range oTicker.C {
			cidToSession.Range(func(sCId string, oSession *atomic.Pointer[Session]) bool {
				oThisSession := oSession.Load()

				oIdle := time.Since(oThisSession.ActivedAt)
				if oIdle <= activeTimeout {
					return true
				}

				oConn := oThisSession.Connection

				pkg.Logger(pkg.WebsocketAdmin).Info(
					"actived timeout, force close",
					zap.String("cid", sCId),
					zap.Duration("idle", oIdle),
				)

				sPointer := fmt.Sprintf("%p", oConn)
				pointerToCid.Del(sPointer)
				cidToSession.Del(sCId)

				cIdsChannels.RemoveLeft(sCId)
				auIdCIds.RemoveRight(sCId)

				oConn.Close()

				return true
			})
		}
	}()

	oMux := http.NewServeMux()
	oMux.Handle("/Admin", oAdminEventer)

	return oMux
}
