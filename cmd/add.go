package cmd

import (

	"github.com/hSATAC/2fa/keychain"
	"github.com/spf13/cobra"
	"os"
	"bufio"
	"log"
	"fmt"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a 2fa account",
	Long: "Add a 2fa account manually",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		account := args[0]

		fmt.Fprintf(os.Stderr, "2fa key for %s: ", account)
		key, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			log.Fatalf("error reading account: %v", err)
		}
		key = key[:len(key)-1] // chop \n
		keychain.Add(account, key)
	},
}

func init() {
	addCmd.SetUsageTemplate(`Usage:

  # Add account manually:
  $ 2fa add [account]
  2fa key for github: [TOTP key]
  `)
	rootCmd.AddCommand(addCmd)
}

