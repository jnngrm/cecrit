package main

import (
    "io/ioutil"
    "github.com/blaskovicz/go-cryptkeeper"
    "encoding/base64"
    "fmt"
    "os"
    "flag"
)

const keyBase64 string = "ffCJ7/JAdIzbsyY+zqIJmyECx5P5LzLKyFepKhzngb0="
var encAction, decAction *bool

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func setKey() {
  key, err := base64.StdEncoding.DecodeString(keyBase64);
  check(err)
  err = cryptkeeper.SetCryptKey([]byte(key))
  check(err)
}

func setup() {
  setKey()
  encAction = flag.Bool("enc", false, "Encrypt files")
  decAction = flag.Bool("dec", false, "Decrypt files")
  flag.Parse()
}

func encryptFile(name string) {
  dat, err := ioutil.ReadFile(name)
  check(err)
  enc, err := cryptkeeper.Encrypt(string(dat))
  check(err)
  encName, err := cryptkeeper.Encrypt(name)
  check(err)
  path := fmt.Sprintf("./enc/%s", encName)
  err = ioutil.WriteFile(path, []byte(enc), 0600)
  check(err)
}

func decryptFile(encName string) {
  name, err := cryptkeeper.Decrypt(encName)
  check(err)
  encPath := fmt.Sprintf("./enc/%s", encName)
  dat, err := ioutil.ReadFile(encPath)
  check(err)
  dec, err := cryptkeeper.Decrypt(string(dat))
  check(err)
  decPath := fmt.Sprintf("./dec/%s", name)
  err = ioutil.WriteFile(decPath, []byte(dec), 0600)
  check(err)
}

func encrypt() {
  _ = os.Mkdir("./enc", 0700)
  files, err := ioutil.ReadDir(".")
  check(err)
  for _, file := range files {
    stat, err := os.Stat(file.Name())
    check(err)
    if stat.Mode().IsRegular() {
      encryptFile(file.Name())
    }
  }
}

func decrypt() {
  _ = os.Mkdir("./dec", 0700)
  files, err := ioutil.ReadDir("./enc")
  check(err)
  for _, file := range files {
    path := fmt.Sprintf("./enc/%s", file.Name())
    stat, err := os.Stat(path)
    check(err)
    if stat.Mode().IsRegular() {
      decryptFile(file.Name())
    }
  }
}

func main() {
    setup()
    if *encAction == true { encrypt() }
    if *decAction == true { decrypt() }
  }
