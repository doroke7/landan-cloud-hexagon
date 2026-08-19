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

// TcpPoolClient 混合版：free-list 交給 sync.Pool（借它 per-P 無鎖快路徑做
// 連線重複利用），硬上限交給 sync.Cond 顧（sync.Pool 本身完全沒有 cap 概念）。
// total 只用來當「同時借出去的連線數」這個名額計數器，不再存實際的閒置連線——
// 閒置連線本體全部交給 sync.Pool 管，get()/put() 都不用再手動操作 slice。
type TcpPoolClient struct {
	pool  sync.Pool // 快路徑：per-P 無鎖 free-list，命中時完全不用碰下面的 mutex
	mutex sync.Mutex
	cond  *sync.Cond
	total int // 目前借出去（get 到 put/discard 之間）的連線數
	max   int // 同時借出去的連線數上限
}

func NewTcpPoolClient() *TcpPoolClient {
	oSelf := &TcpPoolClient{
		max: tcpPoolMaxSize,
	}
	oSelf.cond = sync.NewCond(&oSelf.mutex)
	oSelf.pool.New = func() any {
		return NewTcpConn()
	}
	return oSelf
}

// get 先跟 cond 要一個名額：total 到上限就 cond.Wait() 掛起，
// 被叫醒後重新檢查一次（for 而不是 if，Wait 被叫醒不代表名額一定還在）。
// 名額到手後才問 sync.Pool 要連線——free-list 有現成的就走無鎖快路徑直接拿，
// 沒有的話 sync.Pool 自己呼叫 New（也就是 NewTcpConn）現撥一條。
func (oSelf *TcpPoolClient) get() (*TcpConn, error) {
	oSelf.mutex.Lock()
	for oSelf.total >= oSelf.max {
		oSelf.cond.Wait()
	}
	oSelf.total++
	oSelf.mutex.Unlock()

	oTcpConn, _ := oSelf.pool.Get().(*TcpConn)
	if oTcpConn == nil {
		// New 撥號失敗，名額沒真的用到，要還回去，不然這個名額永久消失
		oSelf.mutex.Lock()
		oSelf.total--
		oSelf.mutex.Unlock()
		oSelf.cond.Signal()
		return nil, ErrTcpDialFailed
	}
	return oTcpConn, nil
}

// put 把還能用的連線丟回 sync.Pool（等下一次 get() 重複利用），
// 並把名額還給 total，叫醒一個可能正在 get() 裡等待的 goroutine。
func (oSelf *TcpPoolClient) put(oTcpConn *TcpConn) {
	oSelf.pool.Put(oTcpConn)

	oSelf.mutex.Lock()
	oSelf.total--
	oSelf.mutex.Unlock()

	oSelf.cond.Signal()
}

// discard 連線壞掉、不能再借出去時呼叫：直接關掉、不丟回 sync.Pool，
// 但名額一樣要還，不然這個名額會永久卡死，池子會越用越小。
func (oSelf *TcpPoolClient) discard(oTcpConn *TcpConn) {
	oTcpConn.Close()

	oSelf.mutex.Lock()
	oSelf.total--
	oSelf.mutex.Unlock()

	oSelf.cond.Signal()
}
