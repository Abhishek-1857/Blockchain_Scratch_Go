package finalizecontroller

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"pop_v1/models"
	"pop_v1/utils"
	"time"

	"github.com/libp2p/go-libp2p"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func Finalize_response() {
	peerId := utils.Client_id
	addr := utils.Client_addr
	var max_true models.Block
	var max_false models.Block
	total := float64(models.VoteCount)
	total = total * 0.70
	var max float64
	for _, value := range models.Groupmap {
		if value.Vote {
			max_true = value.Block
			max++
		} else {
			max_false = value.Block
		}
	}
	var response models.Response
	if max >= total {
		response.Block = max_true
		response.Id = utils.Self_id.String()
		response.Vote = true
	} else {
		var response models.Response
		response.Block = max_false
		response.Vote = false
	}
	clientHost, _ := libp2p.New(libp2p.NoListenAddrs)
	info := peer.AddrInfo{
		ID:    peerId,
		Addrs: []multiaddr.Multiaddr{addr},
	}
	clientHost.Connect(context.Background(), info)
	tr := &http.Transport{}
	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost))

	client := &http.Client{
		Transport: tr,
		Timeout:   20 * time.Millisecond,
	}

	response.Id = utils.Self_id.String()

	requestBody, err := json.Marshal(response)

	if err != nil {
		log.Fatal("Error in serializing the request")
	}

	log.Println("Sending Response to the client !!")
	_, err = client.Post("libp2p://"+info.ID.String()+"/recieveresponse", "application/json", bytes.NewReader(requestBody))
	for err != nil {
		_, err1 := client.Post("libp2p://"+info.ID.String()+"/recieveresponse", "application/json", bytes.NewReader(requestBody))
		err = err1
	}
	clientHost.Close()
	models.Lock.Lock()
	models.T = 0
	models.Groupmap = make(map[string]models.Response)
	models.VoteCount = 0

	models.Lock.Unlock()
}
