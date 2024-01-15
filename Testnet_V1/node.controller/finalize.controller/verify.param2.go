package finalizecontroller

import (
	"pop_v1/blockchain"
	"pop_v1/models"
	"pop_v1/utils"
)

func verify_param2(block models.Block) bool {

	last_block := blockchain.GetlastBlock()
	last_blockhash := utils.GenerateBlockHash(last_block)
	return last_blockhash == block.BlockHeader.Prevhash
}
