package keychain

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/keybase/go-keychain"
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
			fmt.Println("  -",r.Account)
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

	// display like: 1 2 345678
	// for 6~8 digits code
	code = code[:2] + " " + code[2:]
	code = code[:1] + " " + code[1:]
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

func hotp(key []byte, counter uint64, digits int) int {
	h := hmac.New(sha1.New, key)
	binary.Write(h, binary.BigEndian, counter)
	sum := h.Sum(nil)
	v := binary.BigEndian.Uint32(sum[sum[len(sum)-1]&0x0F:]) & 0x7FFFFFFF
	d := uint32(1)
	for i := 0; i < digits && i < 8; i++ {
		d *= 10
	}
	return int(v % d)
}

func totp(key []byte, t time.Time, digits int) int {
	return hotp(key, uint64(t.UnixNano())/30e9, digits)
}

func code(key []byte) string {
	var code int
	code = totp(key, time.Now(), 8)
	return fmt.Sprintf("%0*d", 8, code)
}