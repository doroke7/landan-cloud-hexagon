package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
)

// client 端只用「一條」連線，靠 TcpRequest.RequestId 讓多個並發呼叫共用它。
// 對照 sample/tcp_2（沒有 id、每次都是單一來回）才看得出多路復用多做了什麼，
// 關鍵動作都標了 // IMPORTANT 多路。

type TcpRequest struct {
	RequestId string `json:"request_id"`
	Code      int    `json:"code"`
	Method    string `json:"method"`
	Param     string `json:"param"`
}

type TcpResponse struct {
	RequestId string `json:"request_id"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Result    string `json:"result"`
}

func main() {
	oConn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		panic(err)
	}
	defer oConn.Close()

	oReader := bufio.NewReader(oConn)

	var oWriteMu sync.Mutex

	oPendingMu := sync.Mutex{}
	aPending := make(map[string]chan TcpResponse)

	go func() {
		for {
			var oResp TcpResponse
			if err := decodeFrame(oReader, &oResp); err != nil {
				return
			}

			oPendingMu.Lock()
			oChan, ok := aPending[oResp.RequestId]
			delete(aPending, oResp.RequestId)
			oPendingMu.Unlock()

			if ok {
				oChan <- oResp
			}
			// !ok：id 對不上任何還在等的呼叫方，直接丟掉。
		}
	}()

	/*
	   TCP Client 多路復用其實做了幾個調整
	   1. 呼叫的地方 組裝數據，然後把這個 request-id 註冊一個全局的 數據channel 並且返回 channel pop
	   2. 另外啟動 一個 goruntine 他一直 異步的讀取 資料（讀寫-獨立並行），讀取到黏包後的資料就按照 request-id 寫入對應的 channel 去

	*/

	call := func(sRequestId string, sMethod string, sParam string) (TcpResponse, error) {
		oChan := make(chan TcpResponse, 1)

		oPendingMu.Lock()
		aPending[sRequestId] = oChan
		oPendingMu.Unlock()

		oReq := TcpRequest{
			RequestId: sRequestId,
			Method:    sMethod,
			Param:     sParam,
		}

		oWriteMu.Lock()
		err := encodeFrame(oConn, oReq)
		oWriteMu.Unlock()
		if err != nil {
			return TcpResponse{}, err
		}

		// IMPORTANT 多路：只 blocking 等自己專屬的 channel，不管其他並發中的
		// request 目前處理到哪、也不管 response 實際回來的順序，一定拿到
		// 屬於自己這一筆的結果。
		return <-oChan, nil
	}

	fmt.Println("=== 多路復用測試：同時發出 3 筆 request（id=1,2,3），id=1 的 method 是 Slow ===")
	fmt.Println("=== server 會讓 Slow 睡比較久，response 一定會亂序回來 ===")
	fmt.Println("=== 如果沒用 RequestId 配對，亂序回來的 response 會被誤認成別筆的結果 ===")

	var oWaitGroup sync.WaitGroup
	for i := 1; i <= 3; i++ {
		oWaitGroup.Add(1)
		go func(iId int) {
			defer oWaitGroup.Done()

			sRequestId := strconv.Itoa(iId)
			sMethod := "Fast"
			if iId == 1 {
				sMethod = "Slow"
			}

			oResp, err := call(sRequestId, sMethod, "message-"+sRequestId)
			if err != nil {
				fmt.Println("call failed:", err)
				return
			}

			fmt.Printf("client got response for id=%s: message=%s result=%s\n", sRequestId, oResp.Message, oResp.Result)
		}(i)
	}
	oWaitGroup.Wait()
}

func encodeFrame(oWriter io.Writer, oPayload any) error {
	aBody, err := json.Marshal(oPayload)
	if err != nil {
		return err
	}

	aFrame := make([]byte, 4+len(aBody))
	binary.BigEndian.PutUint32(aFrame[0:4], uint32(len(aBody)))
	copy(aFrame[4:], aBody)

	_, err = oWriter.Write(aFrame)
	return err
}

func decodeFrame(oReader io.Reader, oPayload any) error {
	aLengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(oReader, aLengthBuf); err != nil {
		return err
	}

	iLength := binary.BigEndian.Uint32(aLengthBuf)

	aBody := make([]byte, iLength)
	if _, err := io.ReadFull(oReader, aBody); err != nil {
		return err
	}

	return json.Unmarshal(aBody, oPayload)
}
