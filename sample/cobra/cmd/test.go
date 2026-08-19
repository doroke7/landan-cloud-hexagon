package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
    Use:   "test",
    Short: "Print T$ST message",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("TEST Cobra")
    },
}