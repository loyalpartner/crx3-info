package main

import (
	"crx3-info"
	"log"
)

func main() {
	crx := crx3.NewCrx3("/home/lee/Downloads/3085566c-ab5e-41c3-931b-fff6d49b1146.crx")

	if err := crx.Load(); err != nil {
		log.Fatal(err)
	}

	log.Printf("CrxID: %s", crx.CrxId())
	log.Printf("Magic: %s", crx.Magic)
	log.Printf("Version: %d", crx.Version)
	log.Printf("HeaderSize: %d", crx.HeaderSize)
	log.Printf("Header: %+v", crx.HeaderDetails())
}
