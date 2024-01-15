package accounthelper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"pop_v1/models"
	"pop_v1/utils"
	"pop_v1/wallet"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func AddAccount(w http.ResponseWriter, r *http.Request) {
	publicKey, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var acc_wallet wallet.Wallet
	acc_wallet.PublicKey = publicKey

	pub_addr := fmt.Sprintf("%s", acc_wallet.Address())

	response := models.HttpResponse{}
	response.Data = pub_addr
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	jsonData, err := json.Marshal(acc_wallet)
	go doPublish(context.Background(), utils.AccountTopic, []byte(jsonData))
	time.Sleep(1 * time.Second)

}

func doPublish(ctx context.Context, topic *pubsub.Topic, data []byte) {
	if err := topic.Publish(ctx, data); err != nil {
		log.Println("### Publish error:", err)
	}
}
