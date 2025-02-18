package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "quic",
	Short: "A quick function CLI.",
	Long:  "This CLI helps with quickly exposing helper functions for repitive operations like encoding/decoding a string, generating uuids, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello there")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
