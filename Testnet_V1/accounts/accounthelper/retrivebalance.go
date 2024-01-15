package accounthelper

import (
	"bytes"
	"io"
	"net/http"

	"github.com/libp2p/go-libp2p/core/peer"
)

func Retrievebalance(address string, clientHost *http.Client, client_id peer.ID) (string, error) {

	res, err := clientHost.Post("libp2p://"+client_id.String()+"/getbalancebyaddress", "application/json", bytes.NewReader([]byte(address)))
	for err != nil {
		res1, err1 := clientHost.Post("libp2p://"+client_id.String()+"/getbalancebyaddress", "application/json", bytes.NewReader([]byte(address)))
		res = res1
		err = err1
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err

	}
	return string(bodyBytes), err
}
