package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"pop_v1/models"
)

func GenerateBlockHash(block models.Block) string {
	var blockhash string
	blockheader, _ := json.Marshal(block.BlockHeader)
	data, _ := json.Marshal(block.MetaData)
	blockhash_bytes := sha256.Sum256([]byte(hex.EncodeToString(data) + hex.EncodeToString(blockheader)))
	blockhash = hex.EncodeToString(blockhash_bytes[:])
	return blockhash
}
