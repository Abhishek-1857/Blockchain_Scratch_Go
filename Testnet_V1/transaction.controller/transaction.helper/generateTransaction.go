package transactionhelper

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	Ecdsa "pop_v1/ecdsa"
	"pop_v1/models"

	"github.com/libp2p/go-libp2p/core/peer"
)

func SendTransaction(from string, to string, amount float64, clienthost *http.Client, clienid peer.ID) (string, error) {

	var txn models.TransactionData
	txn.From = from
	txn.To = to
	txn.Amount = amount
	txn.Denom = "blx"
	data, err := json.Marshal(txn)
	if err != nil {
		return "", err
	}
	// signing the transaction
	r, s, err := Ecdsa.SignMessage(data)

	signature := append(r.Bytes(), s.Bytes()...)

	if err != nil {
		return "", err
	}

	var transaction models.Transaction
	transaction.Txn = txn
	transaction.Signature = signature

	transaction_data, err := json.Marshal(transaction)
	if err != nil {
		return "", err
	}

	var transaction_xxxx models.Transaction
	err = json.Unmarshal(transaction_data, &transaction_xxxx)
	if err != nil {
		return "", err
	}
	res, err := clienthost.Post("libp2p://"+clienid.String()+"/sendtesttransaction", "application/json", bytes.NewReader(transaction_data))
	for err != nil {
		res1, err1 := clienthost.Post("libp2p://"+clienid.String()+"/sendtesttransaction", "application/json", bytes.NewReader(transaction_data))
		res = res1
		err = err1
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil

}
