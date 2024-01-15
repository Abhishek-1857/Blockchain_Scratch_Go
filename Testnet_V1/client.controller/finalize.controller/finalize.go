package finalizecontroller

import (
	"context"
	"encoding/json"
	"log"
	"pop_v1/blockchain"
	"pop_v1/models"
	"pop_v1/utils"
	"sync"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var lock sync.Mutex

func Finalize_block() {
	var finalize_hash string = ""
	var max int = 0
	for key, value := range models.VoteMap {
		if len(value) >= max {
			finalize_hash = key
		}
	}
	responses := models.VoteMap[finalize_hash]
	for i := range responses {
		if !responses[i].Vote {
			lock.Lock()
			models.T = 0
			models.TotalResponse = 0
			models.VoteMap = make(map[string][]models.Response)
			lock.Unlock()
			return
		}
	}
	lock.Lock()
	models.T = 0
	models.TotalResponse = 0
	models.VoteMap = make(map[string][]models.Response)
	lock.Unlock()
	var block models.Block
	block.BlockHeader.Height = blockchain.GetlastBlock().BlockHeader.Height + 1
	if len(responses) != 0 {
		block = responses[0].Block
	}

	var jsondata, err = json.Marshal(block)
	if err != nil {
		panic(err)
	}
	go doPublish(context.Background(), utils.BlockchainTopic, []byte(jsondata))
	time.Sleep(1 * time.Second)
}

func doPublish(ctx context.Context, topic *pubsub.Topic, data []byte) {
	if err := topic.Publish(ctx, data); err != nil {
		log.Println("### Publish error:", err)
	}
}
