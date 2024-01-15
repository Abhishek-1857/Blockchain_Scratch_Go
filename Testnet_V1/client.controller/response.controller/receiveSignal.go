package responsecontroller

import (
	"net/http"
	finalizecontroller "pop_v1/client.controller/finalize.controller"
	"pop_v1/models"
	"time"
)

func Sendsignal(http.ResponseWriter, *http.Request) {
	time.Sleep(4 * time.Second)
	lock.Lock()
	models.T = 1
	lock.Unlock()
	go finalizecontroller.Finalize_block()
}
