package models
import "sync"

func GetParamMap() map[string]string {
	var Param_map=map[string]string {
		"0":"Merkleroot",
		"1":"Datahash",
		"2":"Prevhash",	
	}
	return Param_map
}
var Lock sync.Mutex
var T=0
var VoteCount=0
var Groupmap =make(map[string]Response)
var TotalResponse int = 0