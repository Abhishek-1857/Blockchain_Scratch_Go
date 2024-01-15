package blockchain

import (
	"encoding/json"
	"log"
	"pop_v1/database"
	"pop_v1/models"

	"github.com/tecbot/gorocksdb"
)

func Addtoblockchain(block models.Block) {
	data, err := json.Marshal(block)
	if err != nil {
		log.Fatal(err)
	}
	// Writing data to the db
	writeOpts := gorocksdb.NewDefaultWriteOptions()
	defer writeOpts.Destroy()
	err = database.Blockchain_db.Put(writeOpts, []byte("blockid"), data)
	if err != nil {
		log.Fatal("Error writing data: in blockchain-db", err)
	}

	log.Println("Block added successfully to Bockchain")
}
