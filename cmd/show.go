package cmd

import (
	"fmt"
	"github.com/buger/goterm"
	"github.com/hSATAC/2fa/keychain"
	"github.com/pquerna/otp"
	otpTotp "github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
	"log"
	"time"
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
		display(account)
	},
}

func init() {
	showCmd.SetUsageTemplate(`Usage:
  2fa show [account]{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
  `)
	rootCmd.AddCommand(showCmd)
}

func display(account string) {
	url := keychain.Get(account)
	key, err := otp.NewKeyFromURL(url)

	if err != nil {
		log.Fatalf("Error displaying totp when retrieving key: %v", err)
	}
	secret := key.Secret()
	code, err := otpTotp.GenerateCode(secret, time.Now())

	if err != nil {
		log.Fatalf("Error displaying totp when generate code: %v", err)
	}

	period := 30 //Currently, the period parameter is ignored by the Google Authenticator implementations.

	defer func() {
		// Show cursor.
		fmt.Printf("\033[?25h")
	}()
	t := time.Now().UTC().Second()
	countdown := 60 - t
	if countdown > 30 {
		countdown = countdown - 30
	}
	section2 := period / 3 * 2
	section1 := period / 3

	redraw := false

	// TODO: 1. cursor will disappear
	//       2. any key to exit
	for {

		cd := goterm.Bold(fmt.Sprintf("%02d", countdown))
		if countdown <= period && countdown >= section2 {
			cd = goterm.Color(cd, goterm.GREEN)
		} else if countdown < section2 && countdown >= section1 {
			cd = goterm.Color(cd, goterm.YELLOW)
		} else {
			cd = goterm.Color(cd, goterm.RED)
		}

		fmt.Printf("\r      [%s]  %s  \n", cd, code)
		fmt.Printf("\033[?25l")

		time.Sleep(time.Second)

		countdown = countdown - 1
		if countdown < 0 {
			countdown = period
			code, _ = otpTotp.GenerateCode(secret, time.Now())
		}
		redraw = true

		if redraw {
			fmt.Printf("\033[%dA", 1)
		}
	}

}
