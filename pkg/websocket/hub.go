package pkgWebsocket

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
	"go.uber.org/zap"

	pkgUtility "example/pkg/utility"
)

// RootChannel 是連線通過驗證後自動加入的頻道，等同「所有已登入連線」的廣播群。
const RootChannel = "/"

// Session 是一條連線在某個時間點的完整狀態快照。整份是不可變的值——要改任何
// 欄位都是「複製舊快照 → 改副本 → 用 atomic.Pointer 整份換上」，讀的一方
// Load() 到的永遠是某個時間點一致的快照，不會看到改到一半的中間狀態。
type Session struct {
	CId           string
	Connection    *Conn
	ActivedAt     time.Time
	Authenticated bool
	Key           string
	Iv            string
}

// Hub 是單機 websocket 服務的「連線名冊 + 聊天室路由表」，把原本散在各 event
// handler 裡、各自操作四張表的邏輯收斂成一組方法：
//
//   - pointerToCId ：*Conn 的位址字串 → cId。eventer 的 callback 只給 *Conn，
//     要先換回 cId 才能查其他表。
//   - cIdToSession ：cId → *atomic.Pointer[Session]，每條連線一份不可變快照。
//   - auIdCIds     ：登入使用者 id ↔ cId 的多對多關聯（同一個人可以多開）。
//   - cIdChannels  ：cId ↔ 頻道（聊天室房間）的訂閱關係，推播就是查這張表。
//
// 所有方法都可被多個 goroutine 併發呼叫：底層的 hashmap 與 BiMultiMap 各自帶鎖；
// 單一 Session 的更新沿用「單一寫入者」假設——寫的一方都在該連線自己的 read
// loop 裡（OnPong／OnHeartbeat／OnConnect／OnAuthenticate），逾時掃描 goroutine
// 只讀不寫。
type Hub struct {
	pointerToCId *hashmap.Map[string, string]
	cIdToSession *hashmap.Map[string, *atomic.Pointer[Session]]
	auIdCIds     *pkgUtility.BiMultiMap[string, string]
	cIdChannels  *pkgUtility.BiMultiMap[string, string]
}

func NewHub() *Hub {
	return &Hub{
		pointerToCId: hashmap.New[string, string](),
		cIdToSession: hashmap.New[string, *atomic.Pointer[Session]](),
		auIdCIds:     pkgUtility.NewBiMultiMap[string, string](),
		cIdChannels:  pkgUtility.NewBiMultiMap[string, string](),
	}
}

// connKey 用連線指標的位址當 key，等同原本各 handler 裡的 fmt.Sprintf("%p", oConn)。
func connKey(oConn *Conn) string {
	return fmt.Sprintf("%p", oConn)
}

// Add 為剛連上的連線配一個新的 cId 並登記到表上，回傳 cId。極罕見地撞號時
// 回傳 ok=false，呼叫端應直接關掉這條連線。
func (oSelf *Hub) Add(oConn *Conn) (string, bool) {
	sCId := uuid.New().String()

	if _, bExists := oSelf.cIdToSession.Get(sCId); bExists {
		return "", false
	}

	oCell := new(atomic.Pointer[Session])
	oCell.Store(&Session{
		CId:        sCId,
		Connection: oConn,
		ActivedAt:  time.Now(),
	})

	oSelf.cIdToSession.Set(sCId, oCell)
	oSelf.pointerToCId.Set(connKey(oConn), sCId)

	return sCId, true
}

// CId 把 eventer callback 給的 *Conn 換回 cId。
func (oSelf *Hub) CId(oConn *Conn) (string, bool) {
	sKey := connKey(oConn)
	return oSelf.pointerToCId.Get(sKey)
}

// Session 回傳 cId 目前的連線快照。
func (oSelf *Hub) Session(sCId string) (*Session, bool) {
	oCell, bGotten := oSelf.cIdToSession.Get(sCId)
	if !bGotten {
		return nil, false
	}

	oSession := oCell.Load()
	return oSession, true
}

// SessionByConn 是 CId + Session 的組合，是 handler 裡最常見的第一步。
func (oSelf *Hub) SessionByConn(oConn *Conn) (*Session, bool) {
	sCId, bGotten := oSelf.CId(oConn)
	if !bGotten {
		return nil, false
	}

	return oSelf.Session(sCId)
}

// edit 讀出 cId 的舊快照、複製一份交給 fnEdit 改欄位、再用 atomic.Pointer 整份換上。
func (oSelf *Hub) edit(sCId string, fnEdit func(*Session)) {
	oCell, bGotten := oSelf.cIdToSession.Get(sCId)
	if !bGotten {
		return
	}

	oNew := *oCell.Load()
	fnEdit(&oNew)
	oCell.Store(&oNew)
}

// Activate 更新連線的最後活躍時間（OnPong／OnHeartbeat 收到訊息時呼叫）。
func (oSelf *Hub) Activate(sCId string) {
	oSelf.edit(sCId, func(oSession *Session) {
		oSession.ActivedAt = time.Now()
	})
}

// SetKeyIv 存下這條連線之後收發訊息要用的對稱金鑰（OnConnect 帶 K 時）。
func (oSelf *Hub) SetKeyIv(sCId, sKey, sIv string) {
	oSelf.edit(sCId, func(oSession *Session) {
		oSession.Key = sKey
		oSession.Iv = sIv
	})
}

// Authenticate 標記連線通過驗證，把使用者 id 跟 cId 綁起來，並自動加入
// RootChannel。sAuId 是登入後的使用者識別（同一人多開會有多個 cId 綁同一個 id）。
func (oSelf *Hub) Authenticate(sCId, sAuId string) {
	oSelf.edit(sCId, func(oSession *Session) {
		oSession.Authenticated = true
	})

	oSelf.auIdCIds.Insert(sAuId, sCId)
	oSelf.cIdChannels.Insert(sCId, RootChannel)
}

// Authenticated 回報連線是否已通過驗證。
func (oSelf *Hub) Authenticated(sCId string) bool {
	oSession, bGotten := oSelf.Session(sCId)
	return bGotten && oSession.Authenticated
}

// Join 讓連線訂閱一個頻道（聊天室房間）。重複 Join 同一個頻道不會有副作用。
func (oSelf *Hub) Join(sCId, sChannel string) {
	oSelf.cIdChannels.Insert(sCId, sChannel)
}

// Leave 讓連線退出單一頻道。
func (oSelf *Hub) Leave(sCId, sChannel string) {
	oSelf.cIdChannels.Remove(sCId, sChannel)
}

// LeaveAll 讓連線退出目前訂閱的所有頻道（OnUnsubscribe）。
func (oSelf *Hub) LeaveAll(sCId string) {
	oSelf.cIdChannels.RemoveLeft(sCId)
}

// Channels 回傳連線目前訂閱的所有頻道。
func (oSelf *Hub) Channels(sCId string) []string {
	return oSelf.cIdChannels.Left(sCId)
}

// CIds 回傳某頻道目前的所有訂閱者 cId。
func (oSelf *Hub) CIds(sChannel string) []string {
	return oSelf.cIdChannels.Right(sChannel)
}

// AdminUsersInChannel 回傳某頻道裡「不重複」的登入使用者 id（OnPresent 用來列出在場者）。
func (oSelf *Hub) AdminUsersInChannel(sChannel string) []string {
	aMembers := oSelf.CIds(sChannel)

	oSeen := make(map[string]struct{}, len(aMembers))
	aAuIds := make([]string, 0, len(aMembers))

	for _, sCId := range aMembers {
		for _, sAuId := range oSelf.auIdCIds.Right(sCId) {
			if _, bSeen := oSeen[sAuId]; bSeen {
				continue
			}

			oSeen[sAuId] = struct{}{}
			aAuIds = append(aAuIds, sAuId)
		}
	}

	return aAuIds
}

func (oSelf *Hub) AdminUserIds(sAuId string) []string {
	return oSelf.auIdCIds.Left(sAuId)
}

func (oSelf *Hub) Remove(oConn *Conn) string {
	sCId, _ := oSelf.CId(oConn)
	oSelf.drop(sCId, oConn)
	return sCId
}

func (oSelf *Hub) drop(sCId string, oConn *Conn) {
	oSelf.pointerToCId.Del(connKey(oConn))
	oSelf.cIdToSession.Del(sCId)
	oSelf.cIdChannels.RemoveLeft(sCId)
	oSelf.auIdCIds.RemoveRight(sCId)
}

// PublishToCId 把一筆訊息推給單一 cId，連線已不在表上時回傳 false。
func (oSelf *Hub) PublishToCId(sCId string, iType int, aMessage []byte) bool {
	oSession, bGotten := oSelf.Session(sCId)
	if !bGotten {
		return false
	}

	oSession.Connection.WriteMessage(iType, aMessage)
	return true
}

// PublishToChannel 把一筆訊息推給某頻道目前的所有訂閱者（聊天室廣播）。連線已不在
// 表上的自動跳過。
func (oSelf *Hub) PublishToChannel(sChannel string, iType int, aMessage []byte) {
	aMembers := oSelf.CIds(sChannel)
	oSelf.PublishToCIds(aMembers, iType, aMessage)
}

// PublishToAdminUserId 把一筆訊息推給某使用者目前的所有連線（多開時每條都送到）。
func (oSelf *Hub) PublishToAdminUserId(sAuId string, iType int, aMessage []byte) {
	aCIds := oSelf.AdminUserIds(sAuId)
	oSelf.PublishToCIds(aCIds, iType, aMessage)
}

func (oSelf *Hub) PublishToCIds(aCIds []string, iType int, aMessage []byte) {
	for _, sCId := range aCIds {
		oSelf.PublishToCId(sCId, iType, aMessage)
	}
}

// Count 回傳目前在線的連線數。
func (oSelf *Hub) Count() int {
	return oSelf.cIdToSession.Len()
}

// Sweep 掃一輪所有連線，對閒置超過 dTimeout 的強制關閉並清表，回傳被關掉的連線數。
func (oSelf *Hub) Sweep(dTimeout time.Duration) int {
	iClosed := 0

	oSelf.cIdToSession.Range(func(sCId string, oCell *atomic.Pointer[Session]) bool {
		oSession := oCell.Load()

		oIdle := time.Since(oSession.ActivedAt)
		if oIdle <= dTimeout {
			return true
		}

		pkgUtility.Logger(pkgUtility.Websocket).Info(
			"idle timeout, force close",
			zap.String("cid", sCId),
			zap.Duration("idle", oIdle),
		)

		oConn := oSession.Connection
		oSelf.drop(sCId, oConn)
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
