package cmd

import (
	"bufio"
	"fmt"
	"github.com/hSATAC/2fa/keychain"
	"github.com/hSATAC/2fa/screenshot"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

var takeScreenshot bool

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a 2fa account",
	Long: `To add an account, you can either enter the account and
secret manually, or by taking a screenshot of the qrcode.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Basic check
		if takeScreenshot && len(args) > 0 {
			fmt.Println("You can't specify account name when using screenshot.")
			os.Exit(1)
		}

		if !takeScreenshot && len(args) == 0 {
			fmt.Println("You have to enter an account name.")
			os.Exit(1)
		}

		var account, key string
		if takeScreenshot {
			fmt.Println("Please take a screenshot of the qrcode...")
			file, err := ioutil.TempFile(os.TempDir(), "2fa")
			if err != nil {
				log.Fatalf("error creating screenshot file: %v", err)
			}
			defer os.Remove(file.Name())

			// Take screenshot
			err = screenshot.CaptureScreen(file.Name())
			if err != nil {
				log.Fatalf("error taking screenshot: %v", err)
			}

			// Scan screenshot
			account, key = screenshot.ReadQRCode(file.Name())
		} else {
			account := args[0]

			fmt.Fprintf(os.Stderr, "2fa key for %s: ", account)
			key, err := bufio.NewReader(os.Stdin).ReadString('\n')

			if err != nil {
				log.Fatalf("error reading account: %v", err)
			}
			key = key[:len(key)-1] // chop \n
		}

		keychain.Add(account, key)
		fmt.Printf("Account %s has been added.", account)
	},
}

func init() {
	addCmd.Flags().BoolVarP(&takeScreenshot,
		"screenshot", "s", false,
		"Take screenshot of the qrcode to add an account.")
	addCmd.SetUsageTemplate(`Usage:

  # Add account manually:
  $ 2fa add [account]
  2fa key for github: [TOTP key]

  # Add account using screenshot:
  $ 2fa add --screenshot{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
  `)
	rootCmd.AddCommand(addCmd)
}
