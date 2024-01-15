package Ecdsa

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"log"
	"net/http"
	"os"
	accounts "pop_v1/accounts"
	"pop_v1/wallet"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/tecbot/gorocksdb"
)

// StoreAccountsToDB
func save(db *gorocksdb.DB, pubBytes []byte, address accounts.Account) error {
	// Serialize public key

	// Store address JSON in RocksDB with public key as key
	key := pubBytes
	value, err := json.Marshal(address)
	if err != nil {
		return err
	}

	wo := gorocksdb.NewDefaultWriteOptions()
	defer wo.Destroy()

	err = db.Put(wo, key, value)
	if err != nil {
		return err
	}

	return nil
}
func GenerateECDSAKeyPair(clientHost *http.Client, client_id peer.ID) (string, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", err
	}

	publicKey := &privateKey.PublicKey
	err = saveECDSAPrivateKeyToFile("ecdsa_private.pem", privateKey)
	if err != nil {
		return "", err
	}

	var acc_wallet wallet.LocalWallet
	acc_wallet.PrivateKey = *privateKey
	acc_wallet.PublicKey = append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	printECDSAKeyInfo(privateKey, &privateKey.PublicKey)

	res, err := clientHost.Post("libp2p://"+client_id.String()+"/addaccount", "application/json", bytes.NewReader(acc_wallet.PublicKey))
	for err != nil {
		res1, err1 := clientHost.Post("libp2p://"+client_id.String()+"/addaccount", "application/json", bytes.NewReader(acc_wallet.PublicKey))
		res = res1
		err = err1
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(bodyBytes), nil
}

func saveECDSAPrivateKeyToFile(filename string, privateKey *ecdsa.PrivateKey) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	privBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}

	privBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	}

	return pem.Encode(file, privBlock)
}
