package p2p

import (
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// discoveryNotifee handles mDNS peer discovery events
type discoveryNotifee struct {
	peerChan chan peer.AddrInfo
}

// HandlePeerFound is called when a new peer is discovered via mDNS
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	select {
	case n.peerChan <- pi:
	default:
		// Channel is full, skip this peer
	}
}

// InitMDNS initializes mDNS discovery for the given host
func InitMDNS(h host.Host, serviceTag string, peerChan chan peer.AddrInfo) error {
	notifee := &discoveryNotifee{peerChan: peerChan}
	service := mdns.NewMdnsService(h, serviceTag, notifee)
	return service.Start()
}
