package models

import (
	"pop_v1/database"

	"github.com/tecbot/gorocksdb"
)

type TransactionData struct {
	From      string  `json:"from,omitempty"`
	To        string  `json:"to,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
	Denom     string  `json:"denom,omitempty"`
	Nonce     int64   `json:"nonce,omitempty"`
	Txnhash   string  `json:"txnhash,omitempty"`
	Timestamp int64   `json:"timestamp,omitempty"`
	Status    string  `json:"status,omitempty"`
}

type Transaction struct {
	Txn       TransactionData `json:"txn,omitempty"`
	Signature []byte          `json:"signature,omitempty"`
}

func GetLastTransactionID() string {
	readOptions := gorocksdb.NewDefaultReadOptions()
	defer readOptions.Destroy()

	// Create a new iterator and seek to the last key
	iterator := database.TestTransaction_db.NewIterator(readOptions)
	defer iterator.Close()

	iterator.SeekToLast()

	if iterator.Valid() {
		// Extract the last TransactionTest ID and return it
		key := iterator.Key()
		return string(key.Data())
	}

	return ""
}
