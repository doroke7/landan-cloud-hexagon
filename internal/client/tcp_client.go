package client

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"

	bootstrap "example/bootstrap"
	types "example/types"
)

// TcpClient 是最簡單的版本：只有一條連線，沒有連線池（比較 TcpPoolClient），
// 也沒有多路復用（比較 TcpMuxClient）。適合單一 goroutine 循序呼叫的場景；
// 多個 goroutine 同時呼叫會共用同一個 conn/reader，這裡沒有做任何併發保護。

/*
		基本上 這個 tcp client 沒有連結池 也沒有 多路復用，基本上是生產不太能用的，
	   如果高併發情況 多個 代碼請求這個 tcp 是有問題的，可以用鎖保護，但是就會變成 卡頓
*/
type TcpClient struct {
	conn   net.Conn
	reader *bufio.Reader
}

// NewTcpClient 自己讀 bootstrap.CONFIG.CLIENTS.TCP 撥號，不接受外部注入的 net.Conn，
// 呼叫端不用管連線怎麼建的，跟 TcpPoolClient 那種「dial 策略由外部注入」的設計刻意不同。
func NewTcpClient() (*TcpClient, error) {
	oConn, err := net.Dial("tcp", bootstrap.CONFIG.CLIENTS.TCP.HOSTS[0]+":"+bootstrap.CONFIG.CLIENTS.TCP.PORTS[0])
	if err != nil {
		return nil, err
	}

	return &TcpClient{
		conn:   oConn,
		reader: bufio.NewReader(oConn),
	}, nil
}

func (oSelf *TcpClient) Close() error {
	return oSelf.conn.Close()
}

func (oSelf *TcpClient) AdminAuthenticationAuthenticatorSignIn(sName string, sPassword string) (*types.TcpResponse, error) {
	aFrame, err := oSelf.EncodeFrame(types.TcpRequest{
		Method: "Admin.Authentication.Authenticator.SignIn",
		Param:  sName + ":" + sPassword,
	})
	if err != nil {
		log.Println("tcp client: encode failed:", err)
		return nil, err
	}

	if _, err = oSelf.conn.Write(aFrame); err != nil {
		log.Println("tcp client: write failed:", err)
		return nil, err
	}

	var oResp types.TcpResponse
	if err = oSelf.DecodeFrame(oSelf.reader, &oResp); err != nil {
		log.Println("tcp client: decode failed:", err)
		return nil, err
	}

	return &oResp, nil
}

func (oSelf *TcpClient) EncodeFrame(oPayload any) ([]byte, error) {
	oBuf := tcpEncodeBufferPool.Get().(*bytes.Buffer)
	oBuf.Reset()
	defer tcpEncodeBufferPool.Put(oBuf)

	if err := json.NewEncoder(oBuf).Encode(oPayload); err != nil {
		return nil, err
	}

	// json.Encoder.Encode 會多寫一個結尾的 \n，這裡要去掉，行為才會跟 json.Marshal 一致
	aBody := bytes.TrimRight(oBuf.Bytes(), "\n")

	if len(aBody) > tcpMaxBodyLength {
		return nil, ErrTcpBodyTooLarge
	}

	aFrame := make([]byte, 4+len(aBody))
	binary.BigEndian.PutUint32(aFrame[0:4], uint32(len(aBody)))
	copy(aFrame[4:], aBody)

	return aFrame, nil
}

func (oSelf *TcpClient) DecodeFrame(oReader *bufio.Reader, oPayload any) error {
	aLengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(oReader, aLengthBuf); err != nil {
		return err
	}

	iBodyLength := binary.BigEndian.Uint32(aLengthBuf)
	if iBodyLength == 0 || iBodyLength > tcpMaxBodyLength {
		return ErrTcpBodyTooLarge
	}

	pBuf := tcpBodyPool.Get().(*[]byte)
	defer tcpBodyPool.Put(pBuf)

	aBody := (*pBuf)[:iBodyLength]
	if _, err := io.ReadFull(oReader, aBody); err != nil {
		return err
	}

	return json.Unmarshal(aBody, oPayload)
}
