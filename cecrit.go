package main

import (
    "io/ioutil"
    "github.com/blaskovicz/go-cryptkeeper"
    "encoding/base64"
    "fmt"
    "os"
    "flag"
    "crypto/rand"
)

var encAction, decAction *bool
var in, out, keyBase64 *string

func randKey() string {
  keyBytes := make([]byte,32)
  rand.Read(keyBytes)
  return base64.StdEncoding.EncodeToString(keyBytes);
}

func setKey() {
  key, _ := base64.StdEncoding.DecodeString(*keyBase64);
  cryptkeeper.SetCryptKey([]byte(key))
}

func setup() {
  encAction = flag.Bool("enc", false, "Encrypt files")
  decAction = flag.Bool("dec", false, "Decrypt files")
  in = flag.String("in", ".", "Input directory")
  out = flag.String("out", "./out", "Output directory")
  keyBase64 = flag.String("key", randKey(), "32 byte / 256 bit key")
  flag.Parse()
  setKey()
}

func cryptFile(inPath string, name string, outPath string, crypter func(string) (string, error)) {
  filePath := fmt.Sprintf("%s/%s", inPath, name)
  dat, _ := ioutil.ReadFile(filePath)
  content, _ := crypter(string(dat))
  fileName, _ := crypter(name)
  ioutil.WriteFile(fmt.Sprintf("%s/%s", outPath, fileName), []byte(content), 0600)
}

func crypt(inDir string, outDir string, crypter func(string) (string, error)) {
  files, _ := ioutil.ReadDir(inDir)
  for _, file := range files {
    if file.Mode().IsRegular() {
      cryptFile(inDir, file.Name(), outDir, crypter)
    } else if (file.Mode().IsDir()) {
      name, _ := crypter(file.Name())
      os.Mkdir(fmt.Sprintf("%s/%s", outDir, name), 0700)
      crypt(fmt.Sprintf("%s/%s", inDir, file.Name()), fmt.Sprintf("%s/%s", outDir, name), crypter)
    }
  }
}

func main() {
    setup()
    if *encAction == true {
      _ = os.Mkdir(*out, 0700)
      crypt(*in, *out, cryptkeeper.Encrypt)
      fmt.Printf("Key is: %s\n", *keyBase64)
    }
    if *decAction == true {
      _ = os.Mkdir(*out, 0700)
      crypt(*in, *out, cryptkeeper.Decrypt)
    }
  }
