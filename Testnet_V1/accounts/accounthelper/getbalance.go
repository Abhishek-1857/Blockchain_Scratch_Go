package accounthelper

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"pop_v1/models"
	"strconv"
)

func Getbalance(w http.ResponseWriter, r *http.Request) {
	address, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	account, str := Getaccountbyaddress(string(address))
	response := models.HttpResponse{}
	w.Header().Set("Content-Type", "application/json")

	if len(str) > 0 {
		response.Error = "Account not found"

	} else {
		response.Data = strconv.FormatFloat(account.Amount, 'f', -1, 64) + account.Denom
	}
	json.NewEncoder(w).Encode(response)
}
