package cmd

import (
	"github.com/spf13/cobra"
	"github.com/hSATAC/2fa/keychain"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display TOTP code for an account.",
	Long: `Display TOTP code for an account.

The output will look like:

6 8 125305

This is for displaying code in 6-8 digits:
-   125305
-  8125305
- 68125305



`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		account := args[0]
		keychain.Show(account)
	},
}

func init() {
	showCmd.SetUsageTemplate(`Usage:
  2fa show [account]
  `)
	rootCmd.AddCommand(showCmd)
}
