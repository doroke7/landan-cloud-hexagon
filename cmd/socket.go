package cmd

import (
	"github.com/spf13/cobra"

	pkgUtility "example/pkg/utility"
)

var oSocketCommand = &cobra.Command{
	Use:   "socket",
	Short: "啟動 socket 服務",
	Run: func(cmd *cobra.Command, args []string) {
		// 尚未實作：直接 Fatal 讓 exit code 非 0，不要靜默 exit 0。
		oLogger := pkgUtility.Logger(pkgUtility.Socket)
		oLogger.Fatal("socket 服務尚未實作")
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oSocketCommand)
}
