package p2p

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

// Node represents a P2P node in the Waffle network
type Node struct {
	Host     host.Host
	PeerChan chan peer.AddrInfo
	encKey   []byte
}

// NewNode creates a new P2P node with the given configuration
func NewNode(cfg NodeConfig) (*Node, error) {
	// Derive encryption key from private key
	encKey, err := DeriveKey(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive encryption key: %w", err)
	}

	// Build libp2p host options
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", cfg.ListenPort),
		),
	}

	// Create libp2p host
	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p host: %w", err)
	}

	// Create peer channel for discovery
	peerChan := make(chan peer.AddrInfo, 10)

	// Initialize mDNS discovery
	if err := InitMDNS(h, ServiceTag, peerChan); err != nil {
		h.Close()
		return nil, fmt.Errorf("failed to initialize mDNS: %w", err)
	}

	return &Node{
		Host:     h,
		PeerChan: peerChan,
		encKey:   encKey,
	}, nil
}

// DiscoverPeers waits for peers to be discovered within the given timeout
func (n *Node) DiscoverPeers(timeout time.Duration) ([]peer.AddrInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var peers []peer.AddrInfo
	seen := make(map[peer.ID]bool)

	for {
		select {
		case pi := <-n.PeerChan:
			// Skip self
			if pi.ID == n.Host.ID() {
				continue
			}
			// Skip duplicates
			if seen[pi.ID] {
				continue
			}
			seen[pi.ID] = true
			peers = append(peers, pi)
		case <-ctx.Done():
			return peers, nil
		}
	}
}

// EncryptionKey returns the node's encryption key
func (n *Node) EncryptionKey() []byte {
	return n.encKey
}

// Close shuts down the P2P node
func (n *Node) Close() error {
	close(n.PeerChan)
	return n.Host.Close()
}

// ID returns the node's peer ID as a string
func (n *Node) ID() string {
	return n.Host.ID().String()
}

// Addrs returns the node's multiaddresses as strings
func (n *Node) Addrs() []string {
	addrs := n.Host.Addrs()
	result := make([]string, len(addrs))
	for i, addr := range addrs {
		result[i] = addr.String()
	}
	return result
}
