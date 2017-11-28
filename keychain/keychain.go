package keychain

import (
	"encoding/base32"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/keybase/go-keychain"
	otp "github.com/pquerna/otp/totp"
)

const keychainServiceName = "2fa-macOS"

func List() {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(keychainServiceName)
	query.SetAccessGroup(keychainServiceName)
	query.SetMatchLimit(keychain.MatchLimitAll)
	query.SetReturnAttributes(true)
	results, err := keychain.QueryItem(query)

	if err != nil {
		log.Fatal(err)
	}

	if len(results) > 0 {
		fmt.Println("You have following accounts:")
		for _, r := range results {
			fmt.Println("  -", r.Account)
		}
	} else {
		fmt.Println("Run `2fa add` to add acconut.")
	}
}

func Show(account string) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(keychainServiceName)
	query.SetAccount(account)
	query.SetAccessGroup(keychainServiceName)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)
	results, err := keychain.QueryItem(query)
	if err != nil {
		log.Fatalln("keychain query err:", err)
	} else if len(results) != 1 {
		log.Fatalln("keychain query not found:", err)
	}
	code := code(results[0].Data)

	fmt.Printf("%s\n", code)
}

func Add(account string, key string) {
	if _, err := decodeKey(key); err != nil {
		log.Fatalf("invalid key: %v", err)
	}

	label := fmt.Sprintf("%s - %s", keychainServiceName, account)

	item := keychain.NewGenericPassword(
		keychainServiceName,
		account,
		label,
		[]byte(key),
		keychainServiceName)

	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)
	err := keychain.AddItem(item)
	if err != nil {
		log.Fatalf("adding key: %v", err)
	}
}

func decodeKey(key string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(strings.ToUpper(key))
}

func code(key []byte) string {
	code, _ := otp.GenerateCode(string(key), time.Now())
	return code
}
