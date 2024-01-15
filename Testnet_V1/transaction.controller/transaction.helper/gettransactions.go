package transactionhelper

import (
	"encoding/json"
	"net/http"
	"pop_v1/database"
	"pop_v1/models"

	"github.com/tecbot/gorocksdb"
)

func Gettransaction(w http.ResponseWriter, r *http.Request) {
	readOptions := gorocksdb.NewDefaultReadOptions()
	defer readOptions.Destroy()

	iter := database.TestTransaction_db.NewIterator(readOptions)
	defer iter.Close()

	transactions := make([]models.Transaction, 0)

	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
		value := iter.Value()
		transactionData := value.Data()
		var TransactionTest models.Transaction
		if err := json.Unmarshal(transactionData, &TransactionTest); err != nil {
			http.Error(w, "TransactionTest retrival error", http.StatusNotFound)
			return
		}
		transactions = append(transactions, TransactionTest)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
