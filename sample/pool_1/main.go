package main

import (
	"bufio"
	"errors"
	"net"
)

const (
	tcpMaxBodyLength = 1 << 12 // 4KB，避免錯誤/惡意的長度前綴把記憶體打爆
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

// TcpPoolClient 自帶連線池：外部呼叫端不需要知道底層借了哪條連線，呼叫 method 時
// 自動借一條、用完自動還回去。channel 本身就是併發安全的借還機制，多個 goroutine
// 可以同時各自佔用一條連線並發送 request，不會像單一連線+mutex 那樣互相卡住。
type TcpPoolClient struct {
	pool chan *TcpConn
}

// NewTcpPoolClient 建一個空池子，池子大小（最多留幾條閒置連線）讀
// 不會一開始就撥滿。
func NewTcpPoolClient() *TcpPoolClient {
	return &TcpPoolClient{
		pool: make(chan *TcpConn, 24),
	}
}

// get 借一條連線：池子裡有現成的就直接拿，沒有的話（池子空或已達上限被借光）就現撥一條新的，
// 所以真正並發數不受池子大小限制，池子大小只限制「閒置不用時最多留幾條」。
func (oSelf *TcpPoolClient) get() (*TcpConn, error) {
	select {
	case oConn := <-oSelf.pool:
		return oConn, nil
	default:
		return NewTcpConn()
	}
}

// put 把用完的連線還回池子；池子已經滿了就直接把這條連線關掉，不強留超過上限的閒置連線。
func (oSelf *TcpPoolClient) put(oConn *TcpConn) {
	select {
	case oSelf.pool <- oConn:
	default:
		oConn.Close()
	}
}
