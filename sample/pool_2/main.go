package main

import (
	"bufio"
	"errors"
	"net"
	"sync"
)

const (
	tcpMaxBodyLength = 1 << 12 // 4KB，避免錯誤/惡意的長度前綴把記憶體打爆
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

type TcpPoolClient struct {
	pool sync.Pool // sync pool 無法設定最大值，其實不太適合做 連結池
	// sync pool 無法設定最大值，其實不太適合做 連結池
	// sync pool 無法設定最大值，其實不太適合做 連結池
	// sync pool 無法設定最大值，其實不太適合做 連結池
	// sync pool 無法設定最大值，其實不太適合做 連結池
	// sync pool 無法設定最大值，其實不太適合做 連結池
	// sync pool 無法設定最大值，其實不太適合做 連結池
	// sync pool 無法設定最大值，其實不太適合做 連結池
	// sync pool 無法設定最大值，其實不太適合做 連結池
}

func NewTcpPoolClient() *TcpPoolClient {
	return &TcpPoolClient{
		pool: sync.Pool{
			New: func() any {
				return NewTcpConn()
			},
		},
	}
}

func (oSelf *TcpPoolClient) get() (*TcpConn, error) {
	oTcpConn, _ := oSelf.pool.Get().(*TcpConn)
	if oTcpConn == nil {
		return nil, ErrTcpDialFailed
	}
	return oTcpConn, nil
}

// put 把用完的連線還給 sync.Pool；沒有「池子滿了」這種狀態，
// 沒有 Close 的分支——sync.Pool 本身沒有提供「歸還時順便處理掉」的 hook。
func (oSelf *TcpPoolClient) put(oTcpConn *TcpConn) {
	oSelf.pool.Put(oTcpConn)
}
