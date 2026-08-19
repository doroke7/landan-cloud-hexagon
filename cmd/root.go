package cmd

import (
	"os"

	"github.com/spf13/cobra"
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
// 只要該 package 被初始化，這個 package 內所有檔案中的 init() 都會執行一次。
func init() {
	// 在這裡可以定義全局 Flag，例如 --config  ww222
	oRootCommand.PersistentFlags().StringP("config", "c", "config.yaml", "配置文件路徑")
}
