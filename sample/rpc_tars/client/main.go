package main

import (
	"fmt"

	protocolModule "example.rpc.tars/protocol/Module"
	"github.com/TarsCloud/TarsGo/tars"
)

func main() {
	comm := tars.GetCommunicator()

	// 連線字串格式：{app}.{server}.{servant}@{endpoint}，要跟 server 端 conf.xml 裡的
	// app / server / servant 對得上，endpoint 則是直接指定要連的位址（不用另外查註冊中心）
	oHello := protocolModule.NewHelloClient(
		"Module.HelloServer.Hello@tcp -h 127.0.0.1 -p 10015 -t 60000",
		comm,
	)

	var sRsp string
	iRet, err := oHello.TestHello("[I am Tom ]", &sRsp)
	if err != nil {
		panic(err)
	}

	fmt.Println("ret:", iRet, "rsp:", sRsp)
}
