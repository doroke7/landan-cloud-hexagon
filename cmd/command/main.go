package command

import (
	"github.com/spf13/cobra"

	registerCommand "example/internal/register/command"
)

// command 是子命令的父節點，本身不執行動作；不帶 Run，
// 直接下 `command` 會印出子命令列表（跟 root 一樣）。
var Command = &cobra.Command{
	Use:   "command",
	Short: "啟動 Command 命令",
}

func init() {
	Command = registerCommand.Init(Command)
}
