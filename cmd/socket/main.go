package socket

import (
	"github.com/spf13/cobra"

	pkgUtility "example/pkg/utility"
)

var Command = &cobra.Command{
	Use:   "socket",
	Short: "啟動 socket 服務",
	Run: func(cmd *cobra.Command, args []string) {
		// 尚未實作：直接 Fatal 讓 exit code 非 0，不要靜默 exit 0。
		oLogger := pkgUtility.Logger(pkgUtility.Socket)
		oLogger.Fatal("socket 服務尚未實作")
	},
}
