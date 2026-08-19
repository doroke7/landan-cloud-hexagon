package cmd

import (
    "fmt"

	"os"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "A demo CLI app",
    Long:  "This is a demo CLI application using Cobra",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("This is root command")
    },
}

func Execute() {

    if err := rootCmd.Execute(); err != nil {

        fmt.Println(err)

        os.Exit(1)

    }

}