package main

import (
  "flag"
	"fmt"
	"github.com/proglottis/gpgme"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
  "io"
  "path/filepath"
)

const (
TEMPLATE =
`[ACCESS]
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

  } else {
    decrypt()
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

	// close and remove temp file
	tmpfile.Close()
	os.Remove(tmpfile.Name())
}

func decrypt() {
  var filename string
  keys, _ := gpgme.FindKeys("", false)
  for i := 0; i < len(keys); i++ {
    gpgKey := fmt.Sprintf("%s.gpg", keys[i].SubKeys().KeyID())
    filePath , _ := filepath.Abs(gpgKey)
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
  if err != nil {
    panic(err)
  }
  defer plain.Close()
  if _, err := io.Copy(os.Stdout, plain); err != nil {
    panic(err)
  }
}
