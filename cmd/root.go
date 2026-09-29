package cmd

import (
	"os"

	"github.com/spf13/cobra"

	cmdCentrifuge "example/cmd/centrifuge"
	cmdCommand "example/cmd/command"
	cmdCron "example/cmd/cron"
	cmdDaemon "example/cmd/daemon"
	cmdFacade "example/cmd/facade"
	cmdHttp "example/cmd/http"
	cmdRabbitmq "example/cmd/rabbitmq"
	cmdResource "example/cmd/resource"
	cmdSocket "example/cmd/socket"
	cmdSocketio "example/cmd/socketio"
	cmdSource "example/cmd/source"
	cmdTcp "example/cmd/tcp"
	cmdUdp "example/cmd/udp"
	cmdWebsocket "example/cmd/websocket"
)

var oRootCommand = &cobra.Command{
	Use:   "root",
	Short: "高性能後端系統",
}

// Execute 供 main.go 調用
func Execute() {
	if err := oRootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}

// 只要該 package 被初始化，這個 package 內所有檔案中的 init() 都會執行一次。
func init() {
	// 在這裡可以定義全局 Flag，例如 --config  ww222
	oPersistentFlags := oRootCommand.PersistentFlags()
	oPersistentFlags.StringP("config", "c", "config.yaml", "配置文件路徑")

	// 將 socket 指令加入到 root 中
	oRootCommand.AddCommand(cmdSocket.Command)

	// 將 socketio 指令加入到 root 中
	oRootCommand.AddCommand(cmdSocketio.Command)

	// 將 websocket 指令加入到 root 中
	oRootCommand.AddCommand(cmdWebsocket.Command)

	// 將 resource 指令加入到 root 中
	oRootCommand.AddCommand(cmdResource.Command)

	// 將 centrifuge 指令加入到 root 中
	oRootCommand.AddCommand(cmdCentrifuge.Command)

	// 將 cron 指令加入到 root 中
	oRootCommand.AddCommand(cmdCron.Command)

	// 將 daemon 指令加入到 root 中
	oRootCommand.AddCommand(cmdDaemon.Command)

	// 將 facade 指令加入到 root 中
	oRootCommand.AddCommand(cmdFacade.Command)

	// 將 http 指令加入到 root 中
	oRootCommand.AddCommand(cmdHttp.Command)

	// 將 rabbitmq 指令加入到 root 中
	oRootCommand.AddCommand(cmdRabbitmq.Command)

	// 將 source 指令加入到 root 中
	oRootCommand.AddCommand(cmdSource.Command)

	// 將 tcp 指令加入到 root 中
	oRootCommand.AddCommand(cmdTcp.Command)

	// 將 udp 指令加入到 root 中
	oRootCommand.AddCommand(cmdUdp.Command)

	// 將 command 指令加入到 root 中
	oRootCommand.AddCommand(cmdCommand.Command)
}
