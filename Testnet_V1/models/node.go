package models

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// change the struct according to the libp2p addr
type Node struct {
	ID   peer.ID
	Addr multiaddr.Multiaddr
}

type Nodeinfo struct {
	ID   string
	Addr string
}