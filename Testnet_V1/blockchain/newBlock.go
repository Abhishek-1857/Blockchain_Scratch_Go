package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"pop_v1/models"
	"pop_v1/utils"

	"time"
)

func NewBlock(txns []models.Transaction, metadata []string) models.Block {
	// generate the merkle root
	merkleRoot := utils.GenerateMerkleRoot(txns)
	// generate the DataHash
	datahash := utils.GenerateDataHash(txns)
	// find the  previousHash
	PrevHash := GetlastBlock().BlockHeader.Prevhash
	// get the address of current node
	address := utils.Self_addr.String()
	// get the current time
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	//generate the header
	h := GetlastBlock().BlockHeader.Height + 1

	header := models.Header{
		Merkleroot:      merkleRoot,
		Datahash:        datahash,
		Prevhash:        PrevHash,
		Proposeraddress: address,
		Timestamp:       currentTime,
		Height:          h,
		TransactionCnt:  len(metadata),
	}

	var blockhash string
	blockheader, _ := json.Marshal(header)
	data, _ := json.Marshal(metadata)
	blockhash_bytes := sha256.Sum256([]byte(hex.EncodeToString(data) + hex.EncodeToString(blockheader)))
	blockhash = hex.EncodeToString(blockhash_bytes[:])

	block := models.Block{
		BlockHash:   blockhash,
		BlockHeader: header,
		MetaData:    metadata,
	}
	return block
}
