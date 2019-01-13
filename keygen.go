package main

import (
	"encoding/hex"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/proskenion/proskenion/crypto"
	"io/ioutil"
	"log"
)

// https://godoc.org/github.com/jessevdk/go-flags

var opts struct {
	// save to file name
	File string `short:"f" long:"file" description:"A file" value-name:"FILE" default-mask:"-"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	cryptor := crypto.NewEd25519Sha256Cryptor()
	pub, pri := cryptor.NewKeyPairs()
	strPub, strPri := hex.EncodeToString(pub), hex.EncodeToString(pri)

	if opts.File == "" {
		fmt.Println("public_key: 0x" + strPub)
		fmt.Println("private_key: 0x" + strPri)
	} else {
		if err := ioutil.WriteFile(opts.File+".pub", []byte(strPub), 0644); err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile(opts.File+".pri", []byte(strPri), 0644); err != nil {
			log.Fatal(err)
		}
	}
}
