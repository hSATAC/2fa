package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const versionNumber = "v0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of 2fa",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("2fa", versionNumber)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
