package main

import (
	"bufio"
	"fmt"
	"net"

	pkg "example/pkg"
	types "example/types"
)

func main() {
	oConn, err := net.Dial("tcp", "127.0.0.1:4007")
	if err != nil {
		panic(err)
	}
	defer oConn.Close()

	oTcp := pkg.NewTcpRouter()

	aFrame, err := oTcp.EncodeFrame(types.TcpRequest{
		Method: "Admin.Authentication.Authenticator.SignIn",
		Param:  "admin:520999",
	})
	if err != nil {
		panic(err)
	}

	if _, err = oConn.Write(aFrame); err != nil {
		panic(err)
	}

	var oResp types.TcpResponse
	if err := oTcp.DecodeFrame(bufio.NewReader(oConn), &oResp); err != nil {
		panic(err)
	}

	fmt.Printf("code=%d message=%s result=%v\n", oResp.Code, oResp.Message, oResp.Result)
}
