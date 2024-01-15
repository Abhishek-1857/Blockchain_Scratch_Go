package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"pop_v1/database"
	"pop_v1/models"

	"github.com/tecbot/gorocksdb"
)

func Creategenesisblock() {
	// Create a genesis block
	genesisHeader := models.Header{
		Merkleroot:      "",
		Datahash:        "",
		Prevhash:        "",
		Proposeraddress: "",
		Timestamp:       "",
		Height:          1,
		TransactionCnt:  0,
	}

	blockheader, _ := json.Marshal(genesisHeader)
	data, _ := json.Marshal([]string{})
	blockhash_bytes := sha256.Sum256([]byte(hex.EncodeToString(data) + hex.EncodeToString(blockheader)))
	blockhash := hex.EncodeToString(blockhash_bytes[:])
	genesisBlock := models.Block{
		BlockHash:   blockhash,
		BlockHeader: genesisHeader,
		MetaData:    []string{},
	}

	// Convert the genesis block to JSON
	genesisData, err := json.Marshal(genesisBlock)
	if err != nil {
		log.Fatal(err)
	}

	// Writing the genesis block data to the database
	writeOpts := gorocksdb.NewDefaultWriteOptions()
	defer writeOpts.Destroy()
	err = database.Blockchain_db.Put(writeOpts, []byte("abcdefgh"), genesisData)
	if err != nil {
		log.Fatal("Error writing data to blockchain-db", err)
	}

	log.Println("Genesis block added successfully !")
}
