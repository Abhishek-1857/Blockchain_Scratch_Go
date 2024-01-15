package finalizecontroller

import (
	"log"
	"pop_v1/models"
	transactionhelper "pop_v1/transaction.controller/transaction.helper"
	"pop_v1/utils"
)

func verify_param1(block models.Block) bool {

	transaction_list, err := transactionhelper.Retrievetransaction(block.MetaData)
	if err != nil {
		log.Fatal("failed to retrieve transactions")
	}
	generated_merkleroot := utils.GenerateDataHash(transaction_list)
	return block.BlockHeader.Datahash == generated_merkleroot
}
