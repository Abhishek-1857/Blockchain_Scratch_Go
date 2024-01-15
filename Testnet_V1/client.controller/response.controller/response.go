package responsecontroller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	finalizecontroller "pop_v1/client.controller/finalize.controller"
	"pop_v1/models"
	"pop_v1/utils"
	"sync"
	"time"
)

var lock sync.Mutex

func Recieveresponse(w http.ResponseWriter, r *http.Request) {
	var response models.Response
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Fatal(err)
	}

	if models.T == 1 {
		http.Error(w, "message : Timeout", http.StatusNotFound)
		return
	}
	lock.Lock()
	block_recieved := response.Block
	block_hash := utils.GenerateBlockHash(block_recieved)
	response_list := models.VoteMap[block_hash]
	if response_list != nil {
		response_list = append(response_list, response)
		models.VoteMap[block_hash] = response_list

	} else {
		var responses []models.Response
		responses = append(responses, response)
		models.VoteMap[block_hash] = responses
	}
	totalResponses := models.TotalResponse
	models.TotalResponse++
	lock.Unlock()
	if totalResponses == 0 {
		go func() {
			time.Sleep(5 * time.Second)
			lock.Lock()
			models.T = 1
			lock.Unlock()
			go finalizecontroller.Finalize_block()
		}()
	}
}
