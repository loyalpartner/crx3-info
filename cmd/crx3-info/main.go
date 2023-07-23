package main

import (
	"flag"
	"log"

	"github.com/loyalpartner/crx3-info/crx3"
)

var (
	path string
)

func init() {
	flag.StringVar(&path, "path", "", "crx file path")
	flag.Parse()
}

func main() {
	crx, err := crx3.Read(path)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("CrxID: %s", crx.ID)
	log.Printf("Magic: %s", crx.Magic)
	log.Printf("Version: %d", crx.Version)
	log.Printf("HeaderSize: %d", crx.HeaderSize)
	log.Printf("Header: %+v", crx.JsonEncodedHeader())

	if err := crx3.Verify(crx); err != nil {
		log.Printf("%v", err)
	} else {
		log.Printf("Verified: %+v", true)
	}
}
