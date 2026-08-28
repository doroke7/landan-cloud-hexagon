package pkgWebsocket

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"example/bootstrap"
	pkgUtility "example/pkg/utility"
)

// RootChannel 是連線通過驗證後自動加入的頻道，等同「所有已登入連線」的廣播群。
const RootChannel = "/"

// Session 是一條連線在某個時間點的完整狀態快照。整份是不可變的值——要改任何
// 欄位都是「複製舊快照 → 改副本 → 用 atomic.Pointer 整份換上」，讀的一方
// Load() 到的永遠是某個時間點一致的快照，不會看到改到一半的中間狀態。
type Session struct {
	ConnectionId  string
	Connection    *Conn
	ActivedAt     time.Time
	Authenticated bool
	Key           string
	Iv            string
}

// Hub 是單機 websocket 服務的「連線名冊 + 聊天室路由表」，把原本散在各 event
// handler 裡、各自操作四張表的邏輯收斂成一組方法：
//
//   - pointerToConnectionId ：*Conn 的位址字串 → cId。eventer 的 callback 只給 *Conn，
//     要先換回 cId 才能查其他表。
//   - connectionIdToSession ：cId → *atomic.Pointer[Session]，每條連線一份不可變快照。
//   - adminUserIdToConnectionIds     ：登入使用者 id ↔ cId 的多對多關聯（同一個人可以多開）。
//   - connectionIdToChannels  ：cId ↔ 頻道（聊天室房間）的訂閱關係，推播就是查這張表。
//
// 所有方法都可被多個 goroutine 併發呼叫：底層的 hashmap 與 BiMultiMap 各自帶鎖；
// 單一 Session 的更新沿用「單一寫入者」假設——寫的一方都在該連線自己的 read
// loop 裡（OnPong／OnHeartbeat／OnConnect／OnAuthenticate），逾時掃描 goroutine
// 只讀不寫。
type Hub struct {
	pointerToConnectionId      *hashmap.Map[string, string]
	connectionIdToSession      *hashmap.Map[string, *atomic.Pointer[Session]]
	adminUserIdToConnectionIds *pkgUtility.BiMultiMap[string, string]
	connectionIdToChannels     *pkgUtility.BiMultiMap[string, string]
}

func NewHub() *Hub {
	return &Hub{
		pointerToConnectionId:      hashmap.New[string, string](),
		connectionIdToSession:      hashmap.New[string, *atomic.Pointer[Session]](),
		adminUserIdToConnectionIds: pkgUtility.NewBiMultiMap[string, string](),
		connectionIdToChannels:     pkgUtility.NewBiMultiMap[string, string](),
	}
}

// pointer 用連線指標的位址當 key，等同原本各 handler 裡的 fmt.Sprintf("%p", oConn)。
func pointer(oConn *Conn) string {
	return fmt.Sprintf("%p", oConn)
}

// Add 為剛連上的連線配一個新的 cId 並登記到表上，回傳 cId。極罕見地撞號時
// 回傳 ok=false，呼叫端應直接關掉這條連線。
func (oSelf *Hub) Add(oConn *Conn) (string, bool) {
	sConnectionId := uuid.New().String()

	if _, bExists := oSelf.connectionIdToSession.Get(sConnectionId); bExists {
		return "", false
	}

	oCell := new(atomic.Pointer[Session])
	oCell.Store(&Session{
		ConnectionId: sConnectionId,
		Connection:   oConn,
		ActivedAt:    time.Now(),
	})

	oSelf.connectionIdToSession.Set(sConnectionId, oCell)
	oSelf.pointerToConnectionId.Set(pointer(oConn), sConnectionId)

	return sConnectionId, true
}

// ConnectionId 把 eventer callback 給的 *Conn 換回 cId。
func (oSelf *Hub) ConnectionId(oConn *Conn) (string, bool) {
	sKey := pointer(oConn)
	return oSelf.pointerToConnectionId.Get(sKey)
}

// Session 回傳 cId 目前的連線快照。
func (oSelf *Hub) Session(sConnectionId string) (*Session, bool) {
	oCell, bGotten := oSelf.connectionIdToSession.Get(sConnectionId)
	if !bGotten {
		return nil, false
	}

	oSession := oCell.Load()
	return oSession, true
}

// SessionByConn 是 ConnectionId + Session 的組合，是 handler 裡最常見的第一步。
func (oSelf *Hub) SessionByConn(oConn *Conn) (*Session, bool) {
	sConnectionId, bGotten := oSelf.ConnectionId(oConn)
	if !bGotten {
		return nil, false
	}

	return oSelf.Session(sConnectionId)
}

// edit 讀出 cId 的舊快照、複製一份交給 fnEdit 改欄位、再用 atomic.Pointer 整份換上。
func (oSelf *Hub) edit(sConnectionId string, fnEdit func(*Session)) {
	oCell, bGotten := oSelf.connectionIdToSession.Get(sConnectionId)
	if !bGotten {
		return
	}

	oNew := *oCell.Load()
	fnEdit(&oNew)
	oCell.Store(&oNew)
}

// Activate 更新連線的最後活躍時間（OnPong／OnHeartbeat 收到訊息時呼叫）。
func (oSelf *Hub) Activate(sConnectionId string) {
	oSelf.edit(sConnectionId, func(oSession *Session) {
		oSession.ActivedAt = time.Now()
	})
}

// SetKeyIv 存下這條連線之後收發訊息要用的對稱金鑰（OnConnect 帶 K 時）。
func (oSelf *Hub) SetKeyIv(sConnectionId, sKey, sIv string) {
	oSelf.edit(sConnectionId, func(oSession *Session) {
		oSession.Key = sKey
		oSession.Iv = sIv
	})
}

// Authenticate 標記連線通過驗證，把使用者 id 跟 cId 綁起來，並自動加入
// RootChannel。sAdminUserId 是登入後的使用者識別（同一人多開會有多個 cId 綁同一個 id）。
func (oSelf *Hub) Authenticate(sConnectionId, sAdminUserId string) {
	oSelf.edit(sConnectionId, func(oSession *Session) {
		oSession.Authenticated = true
	})

	oSelf.adminUserIdToConnectionIds.Insert(sAdminUserId, sConnectionId)
	oSelf.connectionIdToChannels.Insert(sConnectionId, RootChannel)
}

// Authenticated 回報連線是否已通過驗證。
func (oSelf *Hub) Authenticated(sConnectionId string) bool {
	oSession, bGotten := oSelf.Session(sConnectionId)
	return bGotten && oSession.Authenticated
}

// JoinConnectionIdChannel 讓連線訂閱一個頻道（聊天室房間）。重複 Join 同一個頻道不會有副作用。
func (oSelf *Hub) JoinConnectionIdChannel(sConnectionId, sChannel string) {
	oSelf.connectionIdToChannels.Insert(sConnectionId, sChannel)
}

// LeaveConnectionIdChannel 讓連線退出單一頻道。
func (oSelf *Hub) LeaveConnectionIdChannel(sConnectionId, sChannel string) {
	oSelf.connectionIdToChannels.Remove(sConnectionId, sChannel)
}

// LeaveConnectionId 讓連線退出目前訂閱的所有頻道（OnUnsubscribe）。
func (oSelf *Hub) LeaveConnectionId(sConnectionId string) {
	oSelf.connectionIdToChannels.RemoveLeft(sConnectionId)
}

// ChannelsFromConnectionId 回傳連線目前訂閱的所有頻道。
func (oSelf *Hub) ChannelsFromConnectionId(sConnectionId string) []string {
	return oSelf.connectionIdToChannels.Left(sConnectionId)
}

// ConnectionIdsInChannel 回傳某頻道目前的所有訂閱者 cId。
func (oSelf *Hub) ConnectionIdsInChannel(sChannel string) []string {
	return oSelf.connectionIdToChannels.Right(sChannel)
}

// AdminUsersInChannel 回傳某頻道裡「不重複」的登入使用者 id（OnPresent 用來列出在場者）。
func (oSelf *Hub) AdminUsersInChannel(sChannel string) []string {
	aMembers := oSelf.ConnectionIdsInChannel(sChannel)

	oSeen := make(map[string]struct{}, len(aMembers))
	aAuIds := make([]string, 0, len(aMembers))

	for _, sConnectionId := range aMembers {
		for _, sAdminUserId := range oSelf.adminUserIdToConnectionIds.Right(sConnectionId) {
			if _, bSeen := oSeen[sAdminUserId]; bSeen {
				continue
			}

			oSeen[sAdminUserId] = struct{}{}
			aAuIds = append(aAuIds, sAdminUserId)
		}
	}

	return aAuIds
}

func (oSelf *Hub) AdminUserIds(sAdminUserId string) []string {
	return oSelf.adminUserIdToConnectionIds.Left(sAdminUserId)
}

func (oSelf *Hub) Remove(oConn *Conn) string {
	sConnectionId, _ := oSelf.ConnectionId(oConn)
	oSelf.drop(sConnectionId, oConn)
	return sConnectionId
}

func (oSelf *Hub) drop(sConnectionId string, oConn *Conn) {
	oSelf.pointerToConnectionId.Del(pointer(oConn))
	oSelf.connectionIdToSession.Del(sConnectionId)
	oSelf.connectionIdToChannels.RemoveLeft(sConnectionId)
	oSelf.adminUserIdToConnectionIds.RemoveRight(sConnectionId)
}

// PublishToConnectionId 把一筆訊息推給單一 cId，連線已不在表上時回傳 false。
func (oSelf *Hub) PublishToConnectionId(sConnectionId string, iType int, aByteMessage []byte) bool {
	oSession, bGotten := oSelf.Session(sConnectionId)
	if !bGotten {
		return false
	}

	oSession.Connection.WriteMessage(iType, aByteMessage)
	return true
}

// PublishToChannel 把一筆訊息推給某頻道目前的所有訂閱者（聊天室廣播）。連線已不在
// 表上的自動跳過。
func (oSelf *Hub) PublishToChannel(sChannel string, iType int, aByteMessage []byte) {
	aMembers := oSelf.ConnectionIdsInChannel(sChannel)
	oSelf.PublishToConnectionIds(aMembers, iType, aByteMessage)
}

// PublishToAdminUserId 把一筆訊息推給某使用者目前的所有連線（多開時每條都送到）。
func (oSelf *Hub) PublishToAdminUserId(sAdminUserId string, iType int, aByteMessage []byte) {
	aConnectionIds := oSelf.AdminUserIds(sAdminUserId)
	oSelf.PublishToConnectionIds(aConnectionIds, iType, aByteMessage)
}

func (oSelf *Hub) PublishToConnectionIds(aConnectionIds []string, iType int, aByteMessage []byte) {
	if len(aConnectionIds) <= bootstrap.CONFIG.SERVICES.WEBSOCKET.WORKER {
		for _, sConnectionId := range aConnectionIds {
			oSelf.PublishToConnectionId(sConnectionId, iType, aByteMessage)
		}
		return
	}

	chanWork := make(chan string)
	var oWait sync.WaitGroup

	oWait.Add(bootstrap.CONFIG.SERVICES.WEBSOCKET.WORKER)
	for range bootstrap.CONFIG.SERVICES.WEBSOCKET.WORKER {
		go func() {
			defer oWait.Done()

			for sConnectionId := range chanWork {
				oSelf.PublishToConnectionId(sConnectionId, iType, aByteMessage)
			}
		}()
	}

	for _, sConnectionId := range aConnectionIds {
		chanWork <- sConnectionId
	}
	close(chanWork)

	oWait.Wait()
}

// Count 回傳目前在線的連線數。
func (oSelf *Hub) Count() int {
	return oSelf.connectionIdToSession.Len()
}

// Sweep 掃一輪所有連線，對閒置超過 dTimeout 的強制關閉並清表，回傳被關掉的連線數。
func (oSelf *Hub) Sweep(dTimeout time.Duration) int {
	iClosed := 0

	oSelf.connectionIdToSession.Range(func(sConnectionId string, oCell *atomic.Pointer[Session]) bool {
		oSession := oCell.Load()

		oIdle := time.Since(oSession.ActivedAt)
		if oIdle <= dTimeout {
			return true
		}

		pkgUtility.Logger(pkgUtility.Websocket).Info(
			"idle timeout, force close",
			zap.String("cid", sConnectionId),
			zap.Duration("idle", oIdle),
		)

		oConn := oSession.Connection
		oSelf.drop(sConnectionId, oConn)
		oConn.Close()
		iClosed++

		return true
	})

	return iClosed
}

// SweepCron 每隔 dInterval 掃一次閒置連線，直到 chDone 被關閉。通常在 Init 裡用
// go oHub.SweepCron(...) 起一個背景 goroutine。
func (oSelf *Hub) SweepCron(dInterval, dTimeout time.Duration, chDone <-chan struct{}) {
	oTicker := time.NewTicker(dInterval)
	defer oTicker.Stop()

	for {
		select {
		case <-oTicker.C:
			oSelf.Sweep(dTimeout)
		case <-chDone:
			return
		}
	}
}
