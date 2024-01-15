package utils

import (
	"context"
	"log"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func StreamConsoleTo(ctx context.Context, topic *pubsub.Topic, data []byte) {

	if err := topic.Publish(ctx, data); err != nil {
		log.Println("### Publish error:", err)
	}
}
