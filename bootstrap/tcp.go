package bootstrap

import (
	"net"
	"time"
)

func NewTcp() (net.Conn, error) {
	aHosts := CONFIG.CLIENTS.TCP.HOSTS
	aPorts := CONFIG.CLIENTS.TCP.PORTS

	oDialer := net.Dialer{
		Timeout:   3 * time.Second,  // 建立連線最多等 3 秒
		KeepAlive: 10 * time.Second, // 底層 socket 保活，每 10 秒探測一次，跟 resource client 的 keepalive 同樣目的
	}

	var err error
	for iIndex, sHost := range aHosts {
		sPort := aPorts[0]
		if iIndex < len(aPorts) {
			sPort = aPorts[iIndex]
		}

		var oConn net.Conn
		if oConn, err = oDialer.Dial("tcp", net.JoinHostPort(sHost, sPort)); err == nil {
			return oConn, nil
		}
	}

	return nil, err
}
