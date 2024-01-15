package transactionhelper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"pop_v1/accounts/accounthelper"
	Ecdsa "pop_v1/ecdsa"
	"pop_v1/models"
	"pop_v1/utils"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func Recievetransaction(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	response := models.HttpResponse{}
	w.Header().Set("Content-Type", "application/json")
	var transaction models.Transaction
	err = json.Unmarshal(bodyBytes, &transaction)
	if err != nil {
		response.Error = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	from_acc, str := accounthelper.Getaccountbyaddress(transaction.Txn.From)
	if len(str) > 0 {
		response.Error = "From:" + str
		json.NewEncoder(w).Encode(response)
		return
	}

	_, str1 := accounthelper.Getaccountbyaddress(transaction.Txn.To)
	if len(str1) > 0 {
		response.Error = "To:" + str1
		json.NewEncoder(w).Encode(response)
		return
	}

	if from_acc.Amount < transaction.Txn.Amount {
		response.Error = "Insuffieient Balance"
		json.NewEncoder(w).Encode(response)
		return
	}
	if !from_acc.Txn_pending {

		// transaction data to verify from signature
		txn_data, err := json.Marshal(transaction.Txn)
		if err != nil {
			response.Error = "Serilization Error"
			json.NewEncoder(w).Encode(response)
			return
		}
		// verify the siganture of the transaction
		if !Ecdsa.VerifySignature(from_acc.PubKey, txn_data, transaction.Signature) {
			response.Error = "Transaction failed to verify !!"
			json.NewEncoder(w).Encode(response)
			return
		} else {
			transaction.Txn.Timestamp = time.Now().Unix()
			transaction.Txn.Nonce = from_acc.Nonce + 1
			transaction.Txn.Status = "Pending"
			txn_bytes, _ := json.Marshal(transaction.Txn)
			txn_hash := sha256.Sum256([]byte(txn_bytes))
			transaction.Txn.Txnhash = hex.EncodeToString(txn_hash[:])
			transactionJSON, err := json.Marshal(transaction)
			if err != nil {
				response.Error = "Serialization Error"
				json.NewEncoder(w).Encode(response)
				return
			}

			go doPublish(context.Background(), utils.CTopic, []byte(transactionJSON))
			time.Sleep(1 * time.Second)
			response.Data = transaction.Txn.Txnhash
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		response.Error = "Nonce already used !"
		json.NewEncoder(w).Encode(response)
		return
	}
}

func doPublish(ctx context.Context, topic *pubsub.Topic, data []byte) {
	if err := topic.Publish(ctx, data); err != nil {
		log.Println("### Publish error:", err)
	}
}
