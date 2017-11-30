package cmd

import (
	"fmt"
	"log"
	"time"

	"os/exec"

	"github.com/buger/goterm"
	"github.com/hSATAC/2fa/keychain"
	"github.com/hSATAC/2fa/menu"
	"github.com/pquerna/otp"
	otpTotp "github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display TOTP code for an account.",
	Long: `Display TOTP code for an account.

      [15]  108226

The first 2 digits are the countdown of the TOTP.

Press any key to copy the code and exit.`,
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

	keyEventQueue := make(chan int)

	go func() {
		for {
			ascii, _, _ := menu.GetChar()
			keyEventQueue <- ascii
		}
	}()

	for {
		select {
		case _ = <-keyEventQueue:
			// enter          esc            q
			//if ascii == 13 || ascii == 27 || ascii == 113 {
			//	return
			//}
			err = pbcopy(code)
			if err == nil {
				fmt.Println("\n\n\rTOTP code has been copied.")
			}
			return

		default:
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

}
func pbcopy(text string) error {
	copyCmd := exec.Command("pbcopy")
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return copyCmd.Wait()
}
