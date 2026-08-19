package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Packet 定義我們的 UDP 應用層封包結構
type Packet struct {
	SeqID uint32 // 序號
	Len   uint32 // 資料長度
	Data  []byte // 實際資料
}

// Pack 打包：將結構體轉換為二進位 Byte 陣列
func Pack(seqID uint32, data []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	// 1. 寫入 4 欄位的 SeqID
	err := binary.Write(buf, binary.BigEndian, seqID)
	if err != nil {
		return nil, err
	}

	// 2. 寫入 4 欄位的 Data 長度
	dataLen := uint32(len(data))
	err = binary.Write(buf, binary.BigEndian, dataLen)
	if err != nil {
		return nil, err
	}

	// 3. 寫入實際 Data
	err = binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Unpack 解包：將接收到的 Byte 陣列解析回 Packet 結構
func Unpack(rawData []byte) (*Packet, error) {
	if len(rawData) < 8 {
		return nil, fmt.Errorf("封包長度不足 8 位元組（Header 損壞）")
	}

	buf := bytes.NewReader(rawData)
	packet := &Packet{}

	// 1. 讀取 SeqID
	err := binary.Read(buf, binary.BigEndian, &packet.SeqID)
	if err != nil {
		return nil, err
	}

	// 2. 讀取 Len
	err = binary.Read(buf, binary.BigEndian, &packet.Len)
	if err != nil {
		return nil, err
	}

	// 3. 檢查緩衝區實際長度是否與 Header 標記的長度相符
	// 避免因為網路截斷或封包損壞導致錯誤
	actualDataLen := len(rawData) - 8
	if uint32(actualDataLen) < packet.Len {
		return nil, fmt.Errorf("資料被截斷，預期長度: %d, 實際長度: %d", packet.Len, actualDataLen)
	}

	// 4. 讀取 Data
	packet.Data = make([]byte, packet.Len)
	err = binary.Read(buf, binary.BigEndian, &packet.Data)
	if err != nil {
		return nil, err
	}

	return packet, nil
}
