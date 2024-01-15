package accounthelper

import (
	"encoding/json"
	"log"
	"pop_v1/accounts"
	"pop_v1/database"

	"github.com/tecbot/gorocksdb"
)

func Getaccountbyaddress(address string) (accounts.Account, string) {
	// Create read options
	readOpts := gorocksdb.NewDefaultReadOptions()
	defer readOpts.Destroy()

	// Read data from the database using the ID
	data, err := database.Account_db.Get(readOpts, []byte(address))
	if err != nil {
		log.Fatalf("Error reading data from db:")
	}
	defer data.Free()
	var account accounts.Account
	if data.Size() > 0 {
		//Unmarshal the data into a Account struct
		if err := json.Unmarshal(data.Data(), &account); err != nil {
			log.Fatal("Account failed to unmarshal !!")
		}
		return account, ""
	}
	return account, "Account not found"
}
