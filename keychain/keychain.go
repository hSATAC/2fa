package keychain

import (
	"fmt"
	"log"
	"time"

	macKeychain "github.com/keybase/go-keychain"
	"github.com/pquerna/otp"
	otpTotp "github.com/pquerna/otp/totp"
)

const keychainServiceName = "2fa-macOS"

func List() {
	query := macKeychain.NewItem()
	query.SetSecClass(macKeychain.SecClassGenericPassword)
	query.SetService(keychainServiceName)
	query.SetAccessGroup(keychainServiceName)
	query.SetMatchLimit(macKeychain.MatchLimitAll)
	query.SetReturnAttributes(true)
	results, err := macKeychain.QueryItem(query)

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
	query := macKeychain.NewItem()
	query.SetSecClass(macKeychain.SecClassGenericPassword)
	query.SetService(keychainServiceName)
	query.SetAccount(account)
	query.SetAccessGroup(keychainServiceName)
	query.SetMatchLimit(macKeychain.MatchLimitOne)
	query.SetReturnData(true)
	results, err := macKeychain.QueryItem(query)
	if err != nil {
		log.Fatalln("keychain query err:", err)
	} else if len(results) != 1 {
		log.Fatalln("keychain query not found:", err)
	}
	code := code(results[0].Data)

	fmt.Printf("%s\n", code)
}

func AddByURLString(urlString string) error {
	// check if it's valid url
	key, err := otp.NewKeyFromURL(urlString)
	if err != nil {
		return fmt.Errorf("add OTP URL to keychain error: %v", err)
	}

	account := AccountOfKey(key)

	label := fmt.Sprintf("%s - %s", keychainServiceName, account)

	item := macKeychain.NewGenericPassword(
		keychainServiceName,
		account,
		label,
		[]byte(urlString),
		keychainServiceName)

	item.SetSynchronizable(macKeychain.SynchronizableNo)
	item.SetAccessible(macKeychain.AccessibleWhenUnlocked)
	err = macKeychain.AddItem(item)
	if err != nil {
		log.Fatalf("adding key: %v", err)
	}
	return nil
}

func AccountOfKey(key *otp.Key) string {
	issuer := key.Issuer()
	account := key.AccountName()
	if issuer == "" {
		return account
	} else {
		return fmt.Sprintf("%s - %s", issuer, account)
	}
}

func code(key []byte) string {
	code, _ := otpTotp.GenerateCode(string(key), time.Now())
	return code
}
