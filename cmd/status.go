package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var status = &cobra.Command{
	Use: "status",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Server is Ok ")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(status)
}
