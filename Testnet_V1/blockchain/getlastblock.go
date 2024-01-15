package blockchain

import (
	"encoding/json"
	"log"
	"pop_v1/database"
	"pop_v1/models"

	"github.com/tecbot/gorocksdb"
)

func GetlastBlock() models.Block {
	// Create read options
	readOpts := gorocksdb.NewDefaultReadOptions()
	defer readOpts.Destroy()

	// Create an iterator in reverse order
	iterator := database.Blockchain_db.NewIterator(readOpts)
	defer iterator.Close()

	// Seek to the last key
	iterator.SeekToLast()

	// Check if the iterator is valid
	if iterator.Valid() {
		// Access the data as a byte slice
		blockBytes := iterator.Value().Data()

		// Unmarshal the data into a Block struct
		var lastBlock models.Block
		err := json.Unmarshal(blockBytes, &lastBlock)
		if err != nil {
			log.Fatalf("Error unmarshaling data: of Block")
		}
		return lastBlock
	}

	return models.Block{}
}
