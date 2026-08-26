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
	Connection    *websocket.Conn
	ActivedAt     time.Time
	Authenticated bool
	Key           string
	Iv            string
}

func NewSession(oConn *websocket.Conn) *atomic.Pointer[Session] {
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

		  message    / messaged 當向訊息-不需要ack                    reply ❌            完成
		  chat       / chated
		  broadcast  / broadcasted                                 reply ✅ + broadcast ✅
		  notify    / notified                                     reply ✅

		  refresh     /refreshed                                  reply ✅
		  sub-refresh / sub-refreshed                             reply ✅

--------------------------------------需要檢查是否 authenticated -----------------------------------------------

		  unsubscribe / unsubscribed                              reply ❌           完成

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

	oAdminEventer := pkg.NewWebsocketEventer(websocket.Upgrader{
		CheckOrigin: func(oRequest *http.Request) bool {
			return true
		},
	})

	var (
		pointerToCid = hashmap.New[string, string]()
		cidToSession = hashmap.New[string, *atomic.Pointer[Session]]()

		auIdToCids   = pkg.NewBiMultiMap[string, string]()
		auIdChannels = pkg.NewBiMultiMap[string, string]()
	)

	oAdminEventer.OnOpen(func(oConn *websocket.Conn, iType int) {
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

	oAdminEventer.OnPong(func(oConn *websocket.Conn) {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

		if oSession, bGotten := cidToSession.Get(sCId); bGotten {
			oNew := *oSession.Load()
			oNew.ActivedAt = time.Now()
			oSession.Store(&oNew)
		}
	})

	oAdminEventer.OnConnect(func(oConn *websocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
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

	oAdminEventer.OnHeartbeat(func(oConn *websocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
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
			sCId, _ := pointerToCid.Get(sPointer)

			if oSession, bGotten := cidToSession.Get(sCId); bGotten {
				oNew := *oSession.Load()
				oNew.Authenticated = true
				oSession.Store(&oNew)
			}

			// TODO: 暫時寫死，之後要換成 SignIn 驗證出來的真正 admin_user_id。
			auIdToCids.Insert(sCId, sCId)
		}

		return bOk
	})

	// OnRpc 先把 Method 取出來，但目前還沒有真的依它分派到對應的業務邏輯
	// （resource/AppUser/ShowOne 之類的 Action Route），只回 rpced 確認收到。
	oAdminEventer.OnRpc(func(oConn *websocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
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

	oAdminEventer.OnSubscribe(func(oConn *websocket.Conn, iType int, oWsReq *types.WebsocketRequest) {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

		aAuIds := auIdToCids.Right(sCId)
		if len(aAuIds) == 0 {
			pkg.Logger(pkg.WebsocketAdmin).Info("not authenticated, ignore subscribe", zap.String("cid", sCId))
			return
		}
		sAuId := aAuIds[0]

		var oValue struct {
			Channel string `json:"channel"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json unmarshal error", zap.Error(oErr))
			return
		}

		// 訂閱記在 admin_user_id 這個層級，不是單一 cid——同一個 admin 開好幾個
		// 分頁／裝置都算同一份訂閱，真的要廣播時再透過 cidToSession 反查回實際的
		// 連線。BiMultiMap.Insert 本身有去重、也自己處理並發，不用再額外上鎖。
		auIdChannels.Insert(sAuId, oValue.Channel)

		aByteMessage, oErr := json.Marshal(struct {
			Event   string `json:"event"`
			Channel string `json:"channel"`
			RId     string `json:"r_id"`
		}{
			Event:   "subscribed",
			Channel: oValue.Channel,
			RId:     oWsReq.RId,
		})

		if oErr != nil {
			pkg.Logger(pkg.WebsocketAdmin).Error("json marshal error", zap.Error(oErr))
			return
		}

		oConn.WriteMessage(iType, aByteMessage)
	})

	// OnUnsubscribe 有兩種觸發時機：client 主動送 event: "unsubscribe"，或是連線斷線
	// （iType == -1，跟 OnClose 同一個時間點一起觸發，且會先觸發）。因為沒有帶 oReq，
	// 沒辦法指定「只離開某一個頻道」，兩種情境都當作「這個 admin（sAuId）離開目前
	// 訂閱的全部頻道」——訂閱現在是記在 admin_user_id 這個層級，不是單一 cid。
	// 之後如果要支援「退訂單一頻道」，得替 app event 版本另外設計帶 channel 參數
	// 的簽名。
	oAdminEventer.OnUnsubscribe(func(oConn *websocket.Conn, iType int) {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

		aAuIds := auIdToCids.Right(sCId)
		if len(aAuIds) == 0 {
			return
		}
		sAuId := aAuIds[0]

		// 注意：這裡是整個 auid 一次清掉，沒有做「這個 admin 是不是還有其他 cid
		// 訂閱同一個 channel」的計數——如果同一個 admin 開兩個分頁都訂閱了同一個
		// channel，其中一個分頁退訂／斷線，會連帶把另一個分頁其實還在訂閱的
		// channel 也清掉。之後如果要正確處理多連線共用同一個 auid 的情境，這裡
		// 得改成參照計數。
		auIdChannels.RemoveLeft(sAuId)
	})

	// OnClose 是收不到 cid 的
	oAdminEventer.OnClose(func(oConn *websocket.Conn, iType int) {
		sPointer := fmt.Sprintf("%p", oConn)

		sCId, _ := pointerToCid.Get(sPointer)
		pointerToCid.Del(sPointer)
		cidToSession.Del(sCId)
		auIdToCids.RemoveRight(sCId)

		pkg.Logger(pkg.WebsocketAdmin).Info(
			"disconnected",
			zap.String("cid", sCId),
			zap.Stringer("remoteAddr", oConn.RemoteAddr()),
		)
	})
	oAdminEventer.OnMessage(func(oConn *websocket.Conn, iType int, aMsg []byte) {
		sPointer := fmt.Sprintf("%p", oConn)
		sCId, _ := pointerToCid.Get(sPointer)

		oSession, bGotten := cidToSession.Get(sCId)
		if !bGotten || !oSession.Load().Authenticated {
			pkg.Logger(pkg.WebsocketAdmin).Info("not authenticated, ignore message", zap.Stringer("remoteAddr", oConn.RemoteAddr()))
			return
		}

		// oConn.WriteMessage(iType, aMsg)
	})

	// 定時任務：每 5 分鐘掃一次全部連線，Session.IdleSince() 超過 2 分鐘（沒收到
	// heartbeat 也沒收到 pong）就視為連線壞掉，強制關閉。Close() 之後 ReadMessage 會
	// 出錯、觸發 OnClose，但那是另一個 goroutine 非同步發生的事，這裡不等它，直接把
	// cidToSession 的資料一併刪掉，確保這條 cid 立刻從所有狀態裡消失，不會被下一輪
	// Range 又掃到重複關閉一次。
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
				auIdToCids.RemoveRight(sCId)

				oConn.Close()

				return true
			})
		}
	}()

	oMux := http.NewServeMux()
	oMux.Handle("/Admin", oAdminEventer)

	return oMux
}
