package Ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"math/big"
)

func VerifySignature(publicKey []byte, message []byte, signature []byte) bool {
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:(sigLen / 2)])
	s.SetBytes(signature[(sigLen / 2):])
	// Hash the message
	hash := sha256.Sum256(message)
	curve := elliptic.P256()
	x := big.Int{}
	y := big.Int{}
	keyLen := len(publicKey)
	x.SetBytes(publicKey[:(keyLen / 2)])
	y.SetBytes(publicKey[(keyLen / 2):])
	rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
	// Verify the signature
	return ecdsa.Verify(&rawPubKey, hash[:], &r, &s)
}
