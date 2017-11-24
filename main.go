// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// 2fa is a two-factor authentication agent.
//
// Usage:
//
//	2fa -add [-7] [-8] name
//	2fa -list
//	2fa name
//
// “2fa -add name” adds a new key to the 2fa keychain with the given name.
// It prints a prompt to standard error and reads a two-factor key from standard input.
// Two-factor keys are short case-insensitive strings of letters A-Z and digits 2-7.
//
// The new key generates time-based (TOTP) authentication codes.
//
// By default the new key generates 6-digit codes; the -7 and -8 flags select
// 7- and 8-digit codes instead.
//
// “2fa -list” lists the names of all the keys in the keychain.
//
// “2fa name” prints a two-factor authentication code from the key with the
// given name.
//
// With no arguments, 2fa prints two-factor authentication codes from all
// known time-based keys.
//
// The default time-based authentication codes are derived from a hash of
// the key and the current time, so it is important that the system clock have
// at least one-minute accuracy.
//
// The keychain is stored macOS keychain.
//
// Example
//
// During GitHub 2FA setup, at the “Scan this barcode with your app” step,
// click the “enter this text code instead” link. A window pops up showing
// “your two-factor secret,” a short string of letters and digits.
//
// Add it to 2fa under the name github, typing the secret at the prompt:
//
//	$ 2fa -add github
//	2fa key for github: nzxxiidbebvwk6jb
//	$
//
// Then whenever GitHub prompts for a 2FA code, run 2fa to obtain one:
//
//	$ 2fa github
//	268346
//	$
//
// Or to type less:
//
//	$ 2fa
//	268346	github
//	$
//
package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/keybase/go-keychain"
)

var (
	flagAdd  = flag.Bool("add", false, "add a key")
	flagList = flag.Bool("list", false, "list keys")
	flag7    = flag.Bool("7", false, "generate 7-digit code")
	flag8    = flag.Bool("8", false, "generate 8-digit code")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "\t2fa -add [-7] [-8] keyname\n")
	fmt.Fprintf(os.Stderr, "\t2fa -list\n")
	fmt.Fprintf(os.Stderr, "\t2fa keyname\n")
	os.Exit(2)
}

func main() {
	log.SetPrefix("2fa: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	k := readKeychain()

	if *flagList {
		if flag.NArg() != 0 {
			usage()
		}
		k.list()
		return
	}
	if flag.NArg() == 0 && !*flagAdd {
		usage()
		return
	}
	if flag.NArg() != 1 {
		usage()
	}
	name := flag.Arg(0)
	if strings.IndexFunc(name, unicode.IsSpace) >= 0 {
		log.Fatal("name must not contain spaces")
	}
	if *flagAdd {
		k.add(name)
		return
	}
	k.show(name)
}

type Keychain struct {
	keys map[string]Key
}

type Key struct {
	account string
	raw     []byte
	digits  int
}

const keychainServiceName = "2fa-macOS"

func readKeychain() *Keychain {

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

	c := &Keychain{
		keys: make(map[string]Key),
	}

	for _, r := range results {
		var k Key
		k.account = r.Account
		c.keys[k.account] = k
	}

	return c
}

func (c *Keychain) list() {
	var names []string
	for name := range c.keys {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Println(name)
	}
}

func (c *Keychain) add(name string) {
	size := 6
	if *flag7 {
		size = 7
		if *flag8 {
			log.Fatalf("cannot use -7 and -8 together")
		}
	} else if *flag8 {
		size = 8
	}

	fmt.Fprintf(os.Stderr, "2fa key for %s: ", name)
	text, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalf("error reading key: %v", err)
	}
	text = text[:len(text)-1] // chop \n
	if _, err := decodeKey(text); err != nil {
		log.Fatalf("invalid key: %v", err)
	}

	line := fmt.Sprintf("%s %d %s", name, size, text)

	label := fmt.Sprintf("%s - %s", keychainServiceName, name)

	item := keychain.NewGenericPassword(
		keychainServiceName,
		name,
		label,
		[]byte(line),
		keychainServiceName)

	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)
	err = keychain.AddItem(item)
	if err != nil {
		log.Fatalf("adding key: %v", err)
	}
}

func (c *Keychain) code(name string) string {
	k, ok := c.keys[name]
	if !ok {
		log.Fatalf("no such key %q", name)
	}
	var code int
	code = totp(k.raw, time.Now(), k.digits)
	return fmt.Sprintf("%0*d", k.digits, code)
}

func (c *Keychain) show(name string) {
	if c.keys[name].raw == nil {
		query := keychain.NewItem()
		query.SetSecClass(keychain.SecClassGenericPassword)
		query.SetService(keychainServiceName)
		query.SetAccount(name)
		query.SetAccessGroup(keychainServiceName)
		query.SetMatchLimit(keychain.MatchLimitOne)
		query.SetReturnData(true)
		results, err := keychain.QueryItem(query)
		if err != nil {
			log.Fatalln("keychain query err:", err)
		} else if len(results) != 1 {
			log.Fatalln("keychain query not found:", err)
		} else {
			data := string(results[0].Data)
			datas := strings.Split(data, " ")

			key := Key{}
			key.account = name
			key.digits = int(datas[1][0] - '0')

			raw, err := decodeKey(string(datas[2]))
			if err != nil {
				log.Fatalln("keychain query err:", err)
			}
			key.raw = raw
			c.keys[name] = key
		}
	}
	fmt.Printf("%s\n", c.code(name))
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
