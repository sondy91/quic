package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(jsonCmd)
}

var jsonCmd = &cobra.Command{
	Use:   "json [encode|decode] [string]",
	Short: "A quick function for json related operations",
	Long:  "This command exposes helper functions for json related operations like encoding/decoding a string.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		action := args[0]
		input := args[1]

		switch action {
		case "encode":
			encoded, err := json.Marshal(input)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(encoded))
		case "decode":
			input = strings.TrimSpace(input)
			var decoded interface{}
			err := json.Unmarshal([]byte(input), &decoded)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Decoded JSON: %+v", decoded)
		default:
			fmt.Println("Invalid action. Please use `encode` or `decode`")
		}
	},
}
