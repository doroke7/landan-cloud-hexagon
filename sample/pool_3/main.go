package main

import (
	"bufio"
	"errors"
	"net"
	"sync"
)

const (
	tcpMaxBodyLength = 1 << 12 // 4KB，避免錯誤/惡意的長度前綴把記憶體打爆
	tcpPoolMaxSize   = 24      // 連線總數上限，sync.Cond 版本才做得到「真的設上限」這件事
)

var (
	ErrTcpBodyTooLarge   = errors.New("tcp: body too large")
	ErrTcpMethodNotFound = errors.New("tcp: method not found")
	ErrTcpDialFailed     = errors.New("tcp: dial failed")
)

// tcpConn 是池子裡實際借還的單位：一條 net.Conn 綁一個專屬的 bufio.Reader。
// reader 要跟著連線的生命週期走、不能每次呼叫都重建（理由見 DecodeFrame 註解），
// 所以借出去、還回來都是這一整包一起借還，不會有連線和 reader 對不上的情況。
type TcpConn struct {
	conn   net.Conn
	reader *bufio.Reader
}

// 跟 TcpClient.NewTcpClient 的做法一致。
func NewTcpConn() *TcpConn {
	oConn, err := net.Dial("tcp", "127.0.0.1:8800")
	if err != nil {
		return nil
	}

	return &TcpConn{
		conn:   oConn,
		reader: bufio.NewReader(oConn),
	}
}

func (oSelf *TcpConn) Write(aBuf []byte) (int, error) {
	return oSelf.conn.Write(aBuf)
}

func (oSelf *TcpConn) Close() error {
	return oSelf.conn.Close()
}

// TcpPoolClient 用 sync.Cond 做連接池：跟 sync.Pool 版本最大的差別是，
// 這裡真的可以設定「同時最多存在幾條連線」的硬上限（max）。池子空了、
// 又已經撥到上限時，get() 會用 cond.Wait() 掛起並讓出鎖，直到有人 put()／
// discard() 呼叫 cond.Signal() 才會被叫醒，不會像 sync.Pool 那樣無限現撥新連線。
type TcpPoolClient struct {
	mutex sync.Mutex
	cond  *sync.Cond
	idle  []*TcpConn // 目前閒置、可以直接借出去的連線
	total int        // 目前已經撥出去的連線總數（不管閒置還是借用中），扣掉的時機看 discard
	max   int        // 連線總數上限
}

func NewTcpPoolClient() *TcpPoolClient {
	oSelf := &TcpPoolClient{
		max: tcpPoolMaxSize,
	}
	oSelf.cond = sync.NewCond(&oSelf.mutex)
	return oSelf
}

// get 借一條連線：優先拿現成閒置的；沒有閒置但 total 還沒到上限就現撥一條；
// 兩個條件都不成立（閒置=0 且已經撥滿）就用 cond.Wait() 掛起，
// 被叫醒後從頭重新檢查一次（用 for 而不是 if，因為 Wait 被叫醒不代表條件一定滿足）。
func (oSelf *TcpPoolClient) get() (*TcpConn, error) {
	oSelf.mutex.Lock()
	defer oSelf.mutex.Unlock()

	for {
		if iLen := len(oSelf.idle); iLen > 0 {
			oTcpConn := oSelf.idle[iLen-1]
			oSelf.idle = oSelf.idle[:iLen-1]
			return oTcpConn, nil
		}

		if oSelf.total < oSelf.max {
			oTcpConn := NewTcpConn()
			if oTcpConn == nil {
				return nil, ErrTcpDialFailed
			}
			oSelf.total++
			return oTcpConn, nil
		}

		oSelf.cond.Wait()
	}
}

// put 把還能用的連線還回池子，並叫醒一個可能正在 get() 裡等待的 goroutine。
// 連線壞掉的話要呼叫 discard，不能呼叫這個方法，不然壞連線會被當成閒置連線借出去。
func (oSelf *TcpPoolClient) put(oTcpConn *TcpConn) {
	oSelf.mutex.Lock()
	oSelf.idle = append(oSelf.idle, oTcpConn)
	oSelf.mutex.Unlock()

	oSelf.cond.Signal()
}

// discard 連線壞掉、不能再借出去時呼叫：只把 total 名額還回去讓下一個 get()
// 有機會現撥一條新的頂替，不會因為關掉一條壞連線就讓池子永久少一個名額。
func (oSelf *TcpPoolClient) discard(oTcpConn *TcpConn) {
	oTcpConn.Close()

	oSelf.mutex.Lock()
	oSelf.total--
	oSelf.mutex.Unlock()

	oSelf.cond.Signal()
}
