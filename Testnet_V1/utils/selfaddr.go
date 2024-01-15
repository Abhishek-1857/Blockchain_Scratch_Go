package utils

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

var Self_id peer.ID
var Self_addr multiaddr.Multiaddr
var Client_id peer.ID
var Client_addr multiaddr.Multiaddr
var Serverhost host.Host
var Ps *pubsub.PubSub
