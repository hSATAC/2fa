package cmd

import (
	"github.com/hSATAC/2fa/keychain"
	"github.com/spf13/cobra"
	"github.com/hSATAC/2fa/menu"
	"fmt"
	"os"
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
