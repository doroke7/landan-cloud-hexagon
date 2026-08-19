package main

import (
	protocolModule "example.rpc.tars/protocol/Module"
	service "example.rpc.tars/service" // 這是 tars2go 自動生成的 package
	"github.com/TarsCloud/TarsGo/tars"
	// 這是 tars2go 自動生成的 package
)

func main() {
	cfg := tars.GetServerConfig()

	oProtocolModuleHello := protocolModule.NewHello()

	// 【步驟 C】把「靈魂」注入到「肉體」，並給它一個對外廣播的名字 (Obj)
	sPath := cfg.App + "." + cfg.Server + ".Hello"

	oProtocolModuleHello.AddServantWithContext(service.NewHelloService(), sPath) // <-- 這行就是【注入】！

	// 啟動監聽
	tars.Run()
}
