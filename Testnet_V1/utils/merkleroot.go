package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"pop_v1/models"
)

func GenerateMerkleRoot(transactions []models.Transaction) string {
	if len(transactions) == 0 {
		return ""
	}
	merkleTree := make([]string, len(transactions))

	for i, tx := range transactions {
		// Convert each TransactionTest to bytes
		txBytes := transactionToBytes(tx)

		txHash := sha256.Sum256(txBytes) // Calculate the hash

		// Convert the array to a slice before slicing it
		merkleTree[i] = hex.EncodeToString(txHash[:])
	}

	for len(merkleTree) > 1 {
		var newMerkleTree []string

		for i := 0; i < len(merkleTree); i += 2 {
			left := merkleTree[i]
			var right string
			if i+1 < len(merkleTree) {
				right = merkleTree[i+1]
			} else {
				right = left
			}

			leftHash := sha256.Sum256([]byte(left))
			rightHash := sha256.Sum256([]byte(right))

			combined := hex.EncodeToString(leftHash[:]) + hex.EncodeToString(rightHash[:])

			combinedHash := sha256.Sum256([]byte(combined))

			newMerkleTree = append(newMerkleTree, hex.EncodeToString(combinedHash[:]))
		}

		merkleTree = newMerkleTree
	}

	return merkleTree[0]
}

func transactionToBytes(TransactionTest models.Transaction) []byte {
	data, _ := json.Marshal(TransactionTest)
	return data
}
