package cmd

import (
	"fmt"
	"os"

	"github.com/hSATAC/2fa/keychain"
	"github.com/hSATAC/2fa/menu"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the account names",
	Long: `List all the account names.

Select the account using [up] and [down],
press [Enter] to display the TOTP and the countdown.`,
	Run: func(cmd *cobra.Command, args []string) {
		accounts := keychain.List()
		if len(accounts) == 0 {
			fmt.Println("Run `2fa add` to add an account.")
			os.Exit(0)
		}

		menu := menu.NewButtonMenu("", "Accounts:")
		for _, a := range accounts {
			menu.AddMenuItem(a, a)
		}

		account, escaped := menu.Run()
		if escaped {
			os.Exit(0)
		}

		display(account)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
