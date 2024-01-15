package transactionhelper

import (
	"encoding/json"
	"log"
	"pop_v1/database"
	"pop_v1/models"

	"github.com/tecbot/gorocksdb"
)

func Retrievetransaction(transaction_ids []string) ([]models.Transaction, error) {

	var transaction_list []models.Transaction
	readOptions := gorocksdb.NewDefaultReadOptions()
	defer readOptions.Destroy()

	for i := range transaction_ids {

		data, err := database.TestTransaction_db.Get(readOptions, []byte(transaction_ids[i]))
		if err != nil {
			log.Fatal("transactions can not be retrieved")
			return nil, err
		}

		var TransactionTest models.Transaction
		if err := json.Unmarshal(data.Data(), &TransactionTest); err != nil {
			log.Fatal("transactions can not be retrieved")
		}
		transaction_list = append(transaction_list, TransactionTest)
	}
	return transaction_list, nil
}
