package client

import (
	"bufio"
	"net"
	"strconv"
	"sync"
	"sync/atomic"

	"example/bootstrap"
	pkg "example/pkg"
	types "example/types"
)

// tcpMuxResult 是 readLoop 塞給等待中呼叫方的東西：要嘛是一個完整的 response，
// 要嘛（連線斷掉時）是一個錯誤，兩者只會有一個非空。
type tcpMuxResult struct {
	Resp types.TcpResponse
	Err  error
}

// TcpMuxClient 只用「一條」連線，靠 TcpRequest.RequestId 讓多個並發呼叫共用它，

type TcpMuxClient struct {
	conn   net.Conn
	reader *bufio.Reader
	codec  *pkg.TcpRouter // 只借用它的 EncodeFrame/DecodeFrame，不會拿去 Serve

	writeMu sync.Mutex // 保護「寫一個完整 frame」，多個 goroutine 不能同時寫同一個 socket
	nextId  atomic.Uint64

	pendingMu sync.Mutex
	pending   map[string]chan tcpMuxResult
}

func NewTcpMuxClient() *TcpMuxClient {

	oConn, err := net.Dial("tcp", bootstrap.CONFIG.CLIENTS.TCP.HOSTS[0]+":"+bootstrap.CONFIG.CLIENTS.TCP.PORTS[0])
	if err != nil {
		panic(err)
	}

	oSelf := &TcpMuxClient{
		conn:    oConn,
		reader:  bufio.NewReader(oConn),
		codec:   pkg.NewTcpRouter(),
		pending: make(map[string]chan tcpMuxResult),
	}

	go oSelf.readLoop()

	return oSelf
}

// readLoop 從連線建立到斷線為止只跑這一個 goroutine，其他方法完全不直接碰 oSelf.reader。
func (oSelf *TcpMuxClient) readLoop() {
	for {
		var oResp types.TcpResponse
		if err := oSelf.codec.DecodeFrame(oSelf.reader, &oResp); err != nil {
			oSelf.failAllPending(err)
			return
		}

		oSelf.pendingMu.Lock()
		oChan, ok := oSelf.pending[oResp.RequestId]
		delete(oSelf.pending, oResp.RequestId)
		oSelf.pendingMu.Unlock()

		if ok {
			oChan <- tcpMuxResult{Resp: oResp}
		}
		// !ok：id 對不上任何還在等的呼叫方（例如呼叫方已經因為別的錯誤放棄了），直接丟掉。
	}
}

// failAllPending 在連線壞掉時，讓所有還卡在 call() 裡等待的 goroutine 都能收到錯誤醒過來，
// 而不是永遠 blocking——這是多工設計必須額外處理的部分，單連線單 request 的版本不用管這個。
func (oSelf *TcpMuxClient) failAllPending(err error) {
	oSelf.pendingMu.Lock()
	defer oSelf.pendingMu.Unlock()

	for sId, oChan := range oSelf.pending {
		oChan <- tcpMuxResult{Err: err}
		delete(oSelf.pending, sId)
	}
}

func (oSelf *TcpMuxClient) call(oReq types.TcpRequest) (*types.TcpResponse, error) {
	sId := strconv.FormatUint(oSelf.nextId.Add(1), 10)
	oReq.RequestId = sId

	oChan := make(chan tcpMuxResult, 1)
	oSelf.pendingMu.Lock()
	oSelf.pending[sId] = oChan
	oSelf.pendingMu.Unlock()

	aFrame, err := oSelf.codec.EncodeFrame(oReq)
	if err != nil {
		oSelf.pendingMu.Lock()
		delete(oSelf.pending, sId)
		oSelf.pendingMu.Unlock()
		return nil, err
	}

	oSelf.writeMu.Lock()
	_, err = oSelf.conn.Write(aFrame)
	oSelf.writeMu.Unlock()
	if err != nil {
		oSelf.pendingMu.Lock()
		delete(oSelf.pending, sId)
		oSelf.pendingMu.Unlock()
		return nil, err
	}

	oResult := <-oChan
	if oResult.Err != nil {
		return nil, oResult.Err
	}
	return &oResult.Resp, nil
}

func (oSelf *TcpMuxClient) AdminAuthenticationAuthenticatorSignIn(sName string, sPassword string) (*types.TcpResponse, error) {
	return oSelf.call(types.TcpRequest{
		Method: "Admin.Authentication.Authenticator.SignIn",
		Param:  sName + ":" + sPassword,
	})
}
