package responsecontroller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"pop_v1/models"
)

func Groupresponse(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response models.Response
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Fatal(err)
	}

	models.Lock.Lock()
	models.Groupmap[response.Id] = response
	models.VoteCount++
	models.Lock.Unlock()
}
