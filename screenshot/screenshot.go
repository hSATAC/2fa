package screenshot

import (
	"image/png"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/bieber/barcode"
)

func CaptureScreen(filename string) error {
	cmd := exec.Command("screencapture", "-s", filename)
	err := cmd.Run()
	return err
}

func ReadQRCode(filename string) (account string, key string) {
	fin, err := os.Open(filename)
	defer fin.Close()
	if err != nil {
		log.Fatalf("error reading screenshot file: %v", err)
	}

	src, err := png.Decode(fin)
	if err != nil {
		log.Fatalf("error decoding screenshot file: %v", err)
	}

	img := barcode.NewImage(src)
	scanner := barcode.NewScanner().SetEnabledSymbology(barcode.QRCode, true)

	symbols, err := scanner.ScanImage(img)
	if err != nil {
		log.Fatalf("error scanning qrcode in screenshot file: %v", err)
	}

	if len(symbols) == 0 {
		log.Fatalln("Qrcode not found in screenshot.")
	}

	if len(symbols) > 1 {
		log.Fatalln("Found more than one qrcode in screenshot.")
	}
	return parseOTPAuthURL(symbols[0].Data)
}

func parseOTPAuthURL(otpauthURL string) (account string, key string) {
	if !strings.HasPrefix(otpauthURL, "otpauth://totp/") {
		log.Fatalln("Malformed OTPAuth URL or unsupported HOTP type.")
	}
	u, err := url.Parse(otpauthURL)
	if err != nil {
		log.Fatalln("Malformed OTPAuth URL or unsupported HOTP type.")
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Fatalln("Malformed OTPAuth URL or unsupported HOTP type.")
	}

	// otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example#f
	account = u.Path[1:]
	key = m["secret"][0]

	return account, key
}
