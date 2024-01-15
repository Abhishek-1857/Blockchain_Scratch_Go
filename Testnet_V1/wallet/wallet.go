package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"pop_v1/config"
	"strconv"
	"strings"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

var (
	checkSumlength, _ = strconv.Atoi(config.Config("WALLET_ADDRESS_CHECKSUM"))
	version           = byte(0x00) // hexadecimal representation of zero
)

// https://golang.org/pkg/crypto/ecdsa/
type Wallet struct {
	PublicKey []byte `json:"address,omitempty"`
}
type LocalWallet struct {
	//eliptic curve digital algorithm
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// Validate Wallet Address
func ValidateAddress(address string) bool {

	if len(address) != 34 {
		return false
	}
	//Convert the address to public key hash
	fullHash := Base58Decode([]byte(address))
	// Get the checkSum from Address
	checkSumFromHash := fullHash[len(fullHash)-checkSumlength:]
	//Get the version
	version := fullHash[0]
	pubKeyHash := fullHash[1 : len(fullHash)-checkSumlength]
	checkSum := CheckSum(append([]byte{version}, pubKeyHash...))

	return bytes.Equal(checkSum, checkSumFromHash)
}
func (w *Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)
	versionedHash := append([]byte{version}, pubHash...)
	checksum := CheckSum(versionedHash)
	//version-publickeyHash-checksum
	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	return address
}

// Generate new Key Pair using ecdsa
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub
}

func MakeWallet() *Wallet {
	_, public := NewKeyPair()
	return &Wallet{public}
}
func LocalMakeWallet() *LocalWallet {
	private, public := NewKeyPair()
	return &LocalWallet{private, public}
}
func PublicKeyHash(pubKey []byte) []byte {
	//generate a hash using sha256
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}

	// Re-hash the genrated sha256 using ripemd160
	publicRipMd := hasher.Sum(nil)
	return publicRipMd
}

func CheckSum(data []byte) []byte {
	firstHash := sha256.Sum256(data)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checkSumlength]
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	if err != nil {
		log.Panic(err)
	}
	return decode
}
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}
func PrintWalletAddress(address string, w Wallet) {
	var lines []string
	lines = append(lines, fmt.Sprintf("======ADDRESS:======\n %s ", address))
	lines = append(lines, fmt.Sprintf("======PUBLIC KEY:======\n %x", w.PublicKey))
	log.Println(strings.Join(lines, "\n"))
}
