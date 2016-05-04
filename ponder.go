package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
	"github.com/proglottis/gpgme"
)

// Template for initialized files
const (
	TEMPLATE = `[ACCESS]
%s = *

[myhost]
username = me
password = xx`
	LOCATION = "./"
)

func main() {
	var init bool
	var edit bool

	flag.BoolVar(&init, "i", false, "Initialize a new password db")
	flag.BoolVar(&edit, "e", false, "Edit a password db")

	flag.Parse()

	if init {
		keys, _ := gpgme.FindKeys("", false)
		email := keys[0].UserIDs().Email()

		editString(fmt.Sprintf(TEMPLATE, email))

	} else if edit {
		plain, err := decrypt()
		if err != nil {
			panic(err)
		}
		editString(plain.String())
	} else {
		plain, err := decrypt()
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(os.Stdout, plain); err != nil {
			panic(err)
		}
	}
}

func editString(text string) {

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	tmpfile, _ := ioutil.TempFile("", "")

	_, _ = tmpfile.WriteString(text)

	cmd := exec.Command(editor, tmpfile.Name())

	// without setting std correctly editor will not launch
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	encrypt(tmpfile)
	// close and remove temp file
	tmpfile.Close()
	os.Remove(tmpfile.Name())
}

func findKey(user string, keylist []*gpgme.Key) *gpgme.Key {
	var userKey *gpgme.Key
	for i := 0; i < len(keylist); i++ {
		subkey := keylist[i].SubKeys()
		if user == subkey.KeyID() {
			userKey = keylist[i]
		}
		userIDs := keylist[i].UserIDs()
		if user == userIDs.Email() {
			userKey = keylist[i]
		}
	}
	return userKey
}

func encrypt(tmpFile *os.File) {
	keys, _ := gpgme.FindKeys("", false)

	cfg, err := ini.Load(tmpFile.Read())

	if err != nil {
		panic(err)
	}

	accesshash := cfg.GetSection("ACCESS").KeysHash()

	for user, sections := range accesshash {
		key := findKey(user, keys)
		if key == nil {
			println(fmt.Sprintf("No key found for %s", user))
			continue
		}

		var userSections []string
		if sections != "*" {
			userSections = strings.Split(sections, "")
		}

	}
}

func decrypt() (*bytes.Buffer, error) {
	var filename string
	keys, _ := gpgme.FindKeys("", false)
	for i := 0; i < len(keys); i++ {
		gpgKey := fmt.Sprintf("%s.gpg", keys[i].SubKeys().KeyID())
		filePath, _ := filepath.Abs(gpgKey)
		if _, err := os.Stat(filePath); err == nil {
			filename = filePath
			break
		}
	}

	if filename == "" {
		log.Fatal("Unable to find matching key file")
	}

	f, err := os.Open(filename)
	plain, err := gpgme.Decrypt(f)
	defer plain.Close()
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(plain)
	return buf, err
}
