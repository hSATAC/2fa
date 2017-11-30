package screenshot

import (
	"image/png"
	"log"
	"os"
	"os/exec"

	"github.com/bieber/barcode"
)

func CaptureScreen(filename string) error {
	cmd := exec.Command("screencapture", "-s", filename)
	err := cmd.Run()
	return err
}

func ReadQRCode(filename string) (urlString string) {
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
	return symbols[0].Data
}
