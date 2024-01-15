package noderouter

import (
	"net/http"
	responsecontroller "pop_v1/node.controller/response.controller"
)

func NodeRoutes() {
	http.HandleFunc("/recieve", responsecontroller.Receive)
	http.HandleFunc("/groupresponse", responsecontroller.Groupresponse)
}
