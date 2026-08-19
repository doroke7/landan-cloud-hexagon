package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"time"
)

type Request struct {
	ID      string
	Message string
	SentAt  int64
}

type Response struct {
	ID         string
	Message    string
	ReceivedAt int64
}

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 9000})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("UDP server listening on :9000")

	buf := make([]byte, 64*1024)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		var req Request
		if err := gob.NewDecoder(bytes.NewReader(buf[:n])).Decode(&req); err != nil {
			log.Printf("invalid packet from %s: %v", clientAddr, err)
			continue
		}

		log.Printf("from %s: id=%s message=%q", clientAddr, req.ID, req.Message)

		resp := Response{
			ID:         req.ID,
			Message:    "received: " + req.Message,
			ReceivedAt: time.Now().UnixMilli(),
		}

		var out bytes.Buffer
		if err := gob.NewEncoder(&out).Encode(resp); err != nil {
			log.Println("serialize error:", err)
			continue
		}

		if _, err := conn.WriteToUDP(out.Bytes(), clientAddr); err != nil {
			log.Println("write error:", err)
		}
	}
}
