package main

import (
	"aidanwoods.dev/go-paseto"
	"fmt"
)

func main() {
	secretKey := paseto.NewV4AsymmetricSecretKey() // don't share this!!!
	publicKey := secretKey.Public().ExportHex()    // DO share this one
	privateKey := secretKey.ExportHex()

	fmt.Printf("Public Key: %s\n", publicKey)
	fmt.Printf("Private Key: %s\n", privateKey)
}
