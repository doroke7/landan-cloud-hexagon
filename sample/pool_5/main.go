package main

import (
	"bufio"
	"context"
	"errors"
	"net"
	"sync"

	"golang.org/x/sync/semaphore"
)

const (
	tcpMaxBodyLength = 1 << 12 // 4KB，避免錯誤/惡意的長度前綴把記憶體打爆
	tcpPoolMaxSize   = 24      // 連線總數上限
)

var (
	ErrTcpBodyTooLarge   = errors.New("tcp: body too large")
	ErrTcpMethodNotFound = errors.New("tcp: method not found")
)

// tcpConn 是池子裡實際借還的單位：一條 net.Conn 綁一個專屬的 bufio.Reader。
// reader 要跟著連線的生命週期走、不能每次呼叫都重建（理由見 DecodeFrame 註解），
// 所以借出去、還回來都是這一整包一起借還，不會有連線和 reader 對不上的情況。
type TcpConn struct {
	conn   net.Conn
	reader *bufio.Reader
}

// 跟 TcpClient.NewTcpClient 的做法一致。
func NewTcpConn() (*TcpConn, error) {
	oConn, err := net.Dial("tcp", "127.0.0.1:8800")
	if err != nil {
		return nil, err
	}

	return &TcpConn{
		conn:   oConn,
		reader: bufio.NewReader(oConn),
	}, nil
}

func (oSelf *TcpConn) Write(aBuf []byte) (int, error) {
	return oSelf.conn.Write(aBuf)
}

func (oSelf *TcpConn) Close() error {
	return oSelf.conn.Close()
}

// TcpPoolClient 用 golang.org/x/sync/semaphore 做連接池：用信號量限制「同時最多
// 幾條連線在使用中」（tcpPoolMaxSize），概念上跟 sync.Cond 版本的 max 是同一件事，
// 差別是等待/喚醒交給 semaphore.Weighted 處理，不用自己手寫 cond.Wait()/Signal()，
// 而且 Acquire 吃 context，天生就支援「最多等信號量多久」或「外部取消就不等了」。
// 信號量只負責「准不准借」，「借到的是哪一條連線」還是要自己用 mutex+idle 維護。
type TcpPoolClient struct {
	sem   *semaphore.Weighted
	mutex sync.Mutex
	idle  []*TcpConn
}

func NewTcpPoolClient() *TcpPoolClient {
	return &TcpPoolClient{
		sem: semaphore.NewWeighted(tcpPoolMaxSize),
	}
}

// get 先跟信號量要一個名額：池子滿了就卡在 Acquire 裡，直到有人 put()／discard()
// 釋放名額，或是 ctx 被取消/超時。拿到名額後優先用現成閒置的連線，沒有才現撥新的。
func (oSelf *TcpPoolClient) get(ctx context.Context) (*TcpConn, error) {
	if err := oSelf.sem.Acquire(ctx, 1); err != nil {
		return nil, err
	}

	oSelf.mutex.Lock()
	if iLen := len(oSelf.idle); iLen > 0 {
		oTcpConn := oSelf.idle[iLen-1]
		oSelf.idle = oSelf.idle[:iLen-1]
		oSelf.mutex.Unlock()
		return oTcpConn, nil
	}
	oSelf.mutex.Unlock()

	oTcpConn, err := NewTcpConn()
	if err != nil {
		oSelf.sem.Release(1) // 撥號失敗，名額沒真的用到，要還回去，不然永久少一個名額
		return nil, err
	}
	return oTcpConn, nil
}

// put 把還能用的連線放回閒置池，並釋放一個信號量名額給下一個等待中的 get()。
func (oSelf *TcpPoolClient) put(oTcpConn *TcpConn) {
	oSelf.mutex.Lock()
	oSelf.idle = append(oSelf.idle, oTcpConn)
	oSelf.mutex.Unlock()

	oSelf.sem.Release(1)
}

// discard 連線壞掉時呼叫：關掉連線、不放回 idle，但信號量名額一樣要釋放，
// 不然這個名額會永久卡死，池子會越用越小。
func (oSelf *TcpPoolClient) discard(oTcpConn *TcpConn) {
	oTcpConn.Close()
	oSelf.sem.Release(1)
}
