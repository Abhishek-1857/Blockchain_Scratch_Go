package finalizecontroller

import "pop_v1/models"

func CheckParameter(param_id string, block models.Block) bool {
	my_paramap := models.GetParamMap()
	my_param := my_paramap[param_id]

	if my_param == "Merkleroot" {
		if verify_param0(block) {
			return true
		}
	} else if my_param == "Datahash" {
		if verify_param1(block) {
			return true
		}

	} else if my_param == "Prevhash" {
		if verify_param1(block) {
			return true
		}
	}
	return false
}
