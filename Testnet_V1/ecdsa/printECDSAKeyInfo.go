package Ecdsa

import (
	"crypto/ecdsa"
	"log"
)

func printECDSAKeyInfo(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) {

	log.Printf("X (public key x-coordinate): %x\n", publicKey.X)
	log.Printf("Y (public key y-coordinate): %x\n", publicKey.Y)
	log.Printf("public key hash : %x\n", append(publicKey.X.Bytes(), publicKey.Y.Bytes()...))
}
