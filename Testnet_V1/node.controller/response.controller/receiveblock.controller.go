package responsecontroller

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"pop_v1/models"
	finalizecontroller "pop_v1/node.controller/finalize.controller"
	"pop_v1/utils"
	"time"

	"github.com/libp2p/go-libp2p"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func Receive(w http.ResponseWriter, r *http.Request) {
	log.Println("Request Came to nodes !!")
	if models.T == 1 {
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var request models.Superior_res

	err = json.Unmarshal(bodyBytes, &request)
	if err != nil {
		log.Fatal(err)
	}
	block := request.Block
	recieved_block_hash := utils.GenerateBlockHash(block)
	var response models.Response
	response.Block = block
	response.Id = utils.Self_id.String()
	if recieved_block_hash != request.Block_hash {
		response.Vote = false
	} else {
		if finalizecontroller.CheckParameter(request.Pid, block) {
			response.Vote = true
		} else {
			response.Vote = false
		}

	}
	clientHost, _ := libp2p.New(libp2p.NoListenAddrs)
	// orginal admins id and addr
	addr, _ := multiaddr.NewMultiaddr(request.A_addr)
	peerId, _ := peer.Decode(request.Aid)
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
	requestBody, err := json.Marshal(response)
	if err != nil {
		log.Fatal("message : resposne can not be serilaized")
	}
	// node itselif is the admin

	if request.Aid == utils.Self_id.String() {
		models.Lock.Lock()
		models.Groupmap[utils.Self_id.String()] = response
		models.VoteCount++
		models.Lock.Unlock()
		time.Sleep(2 * time.Second)
		models.Lock.Lock()
		models.T = 1
		models.Lock.Unlock()
		finalizecontroller.Finalize_response()
		return
	}
	_, err = client.Post("libp2p://"+info.ID.String()+"/groupresponse", "application/json", bytes.NewReader(requestBody))
	for err != nil {
		_, err1 := client.Post("libp2p://"+info.ID.String()+"/groupresponse", "application/json", bytes.NewReader(requestBody))
		err = err1
	}
	clientHost.Close()
}
