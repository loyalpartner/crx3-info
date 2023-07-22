package main

import (
	"flag"
	"log"

	"github.com/loyalpartner/crx3-info"
)

var (
	path string
)

func main() {

	flag.StringVar(&path, "path", "", "crx file path")
	flag.Parse()

	crx := crx3.NewCrx3(path)

	if err := crx.Load(); err != nil {
		log.Fatal(err)
	}

	log.Printf("CrxID: %s", crx.CrxId())
	log.Printf("Magic: %s", crx.Magic)
	log.Printf("Version: %d", crx.Version)
	log.Printf("HeaderSize: %d", crx.HeaderSize)
	log.Printf("Header: %+v", crx.JsonEncodedHeader())

	err := crx.Verify()
	if err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("Verified: %+v", true)
	}
}
