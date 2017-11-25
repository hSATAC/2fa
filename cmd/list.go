package cmd

import (
	"github.com/hSATAC/2fa/keychain"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the account names",
	Long: `List all the account names.

If you haven't add any accounts, run:

2fa add

to add an account.`,
	Run: func(cmd *cobra.Command, args []string) {
		keychain.List()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
