package transactionhelper

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"pop_v1/blockchain"
	"pop_v1/config"
	"pop_v1/models"
	superiorhelper "pop_v1/superior.controller/helper"
	"pop_v1/utils"
	"strconv"
	"time"

	"github.com/libp2p/go-libp2p"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func TransactionLogic() {
	//// take the transactions from the memepool
	log.Println("BlockCreation Request Received !!")
	out, done, wg := utils.GetMempoolTxns()
	var txnsMetaData []string
	var txns []models.Transaction
	err := -1

	go func() {
		for {
			select {
			case res := <-done:
				if !res {
					//error is there
					err = 1
				} else {
					err = 0
				}
				return
			case fun := <-out:
				metadata, txn := fun()
				txns = append(txns, txn)
				txnsMetaData = append(txnsMetaData, metadata)
				(*wg).Done()
			}
		}
	}()
	ips := utils.ShowPeers(utils.Ps, utils.Serverhost)
	//// make sure ips and txns are collected.
	for {
		if err > -1 {
			break
		}
	}
	if err == 1 {
		log.Fatal("message : Internal Server Error")
	}

	//// Generate Groups
	superiorhelper.GroupGenerator(&ips)
	groups, _ := strconv.Atoi(config.Config("GROUPS"))
	param, _ := strconv.Atoi(config.Config("PARAM"))

	//// select admins
	var admins []int
	for i := 0; i < param; i++ {
		random := rand.Intn(100000000) % groups
		admins = append(admins, random)
	}
	//// create the block
	block := blockchain.NewBlock(txns, txnsMetaData)
	hash := utils.GenerateBlockHash(block)

	//// brodcast the block
	for i := 0; i < groups; i++ {
		for j := 0; j < param; j++ {
			//i-th group and j-th param
			//for the jth parameter admins[j] is the admin
			go func(I int, J int) {
				ip := ips[I*param+J]
				clientHost, _ := libp2p.New(libp2p.NoListenAddrs)
				peer_id, _ := peer.Decode(ip.ID)
				peer_addr, _ := multiaddr.NewMultiaddr(ip.Addr)
				info := peer.AddrInfo{
					ID:    peer_id,
					Addrs: []multiaddr.Multiaddr{peer_addr},
				}

				clientHost.Connect(context.Background(), info)
				tr := &http.Transport{}
				tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost))
				client := &http.Client{
					Transport: tr,
					Timeout:   60 * time.Millisecond,
				}
				var request models.Superior_res
				request.Aid = ips[admins[J]*param+J].ID
				request.A_addr = ips[admins[J]*param+J].Addr
				request.Gid = strconv.Itoa(I)
				request.Pid = strconv.Itoa(J)
				request.Block_hash = hash
				request.Block = block

				requestBody, Err := json.Marshal(request)
				if Err != nil {
					log.Fatal("Error marshaling Data: request", err)
				}

				_, err := client.Post("libp2p://"+info.ID.String()+"/recieve", "application/json", bytes.NewReader(requestBody))
				for err != nil {
					_, err1 := client.Post("libp2p://"+info.ID.String()+"/recieve", "application/json", bytes.NewReader(requestBody))
					err = err1
				}
				clientHost.Close()
			}(i, j)
		}
	}
}
