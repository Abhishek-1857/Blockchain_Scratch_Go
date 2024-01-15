package utils

import (
	"pop_v1/config"
	"pop_v1/models"
	"strconv"
	"strings"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

func ShowPeers(ps *pubsub.PubSub, node host.Host) []models.Nodeinfo {
	var ips []models.Nodeinfo
	param, _ := strconv.Atoi(config.Config("PARAM"))
	gps, _ := strconv.Atoi(config.Config("GROUPS"))
	CNT := param * gps
	const TopicName = "BlockMagix"
	peers := ps.ListPeers(TopicName)
	for _, p := range peers {
		if Client_id != p {
			multiaddrs := node.Peerstore().Addrs(p)
			for _, addr := range multiaddrs {
				if !strings.Contains(addr.String(), "127.0.0.1") {
					var ip = models.Nodeinfo{
						ID:   p.String(),
						Addr: addr.String(),
					}
					ips = append(ips, ip)
				}
			}
		}
		if len(ips) == CNT {
			break
		}
	}
	return ips
}
