package transactionhelper

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"pop_v1/database"
	"pop_v1/models"

	"github.com/tecbot/gorocksdb"
)

func Gettransactionbyhash(w http.ResponseWriter, r *http.Request) {
	txn_hash, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Create read options
	readOpts := gorocksdb.NewDefaultReadOptions()
	defer readOpts.Destroy()

	// Read data from the database using the ID
	data, err := database.TestTransaction_db.Get(readOpts, txn_hash)
	if err != nil {
		log.Fatalf("Error reading data from db:")
	}

	defer data.Free()
	//Unmarshal the data into a Account struct
	var transaction models.Transaction
	// response := make(map[string]string)
	response := models.HttpResponse{}
	w.Header().Set("Content-Type", "application/json")

	if data.Size() > 0 {
		if err := json.Unmarshal(data.Data(), &transaction); err != nil {
			log.Fatal("Failed to marshal !")
		}
		response.Data = transaction.Txn
	} else {
		response.Error = "transaction hash not found "
	}
	json.NewEncoder(w).Encode(response)
}
