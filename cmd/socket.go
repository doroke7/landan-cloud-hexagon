package cmd

import (
	"github.com/spf13/cobra"
)

var oSocketCommand = &cobra.Command{
	Use:   "socket",
	Short: "啟動 socket 服務",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oSocketCommand)
}
