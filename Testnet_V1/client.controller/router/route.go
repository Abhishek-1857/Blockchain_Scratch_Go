package clientrouter

import (
	"net/http"
	accounthelper "pop_v1/accounts/accounthelper"
	responsecontroller "pop_v1/client.controller/response.controller"
)

func SetupClientRoute() {
	http.HandleFunc("/recieveresponse", responsecontroller.Recieveresponse)
	http.HandleFunc("/signal", responsecontroller.Sendsignal)
	http.HandleFunc("/addaccount", accounthelper.AddAccount)
}
