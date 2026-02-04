package p2p

import "github.com/libp2p/go-libp2p/core/protocol"

// Service discovery constants
const (
	ServiceTag = "waffle-p2p-service"
	ProtocolID = protocol.ID("/waffle/v1")
)

// SyrupPayload represents the encrypted payload structure for P2P communication
type SyrupPayload struct {
	Prompt          string `json:"prompt,omitempty"`
	Data            []byte `json:"data"`
	RequestID       uint64 `json:"request_id,omitempty"`       // Blockchain request ID
	ProviderAddress string `json:"provider_address,omitempty"` // Provider's wallet address
	TokenUsage      uint64 `json:"token_usage,omitempty"`      // AI token usage
	Error           string `json:"error,omitempty"`            // Error message if any
}

// NodeConfig holds configuration for creating a new P2P node
type NodeConfig struct {
	PrivateKey string // Hex-encoded private key for encryption key derivation
	ListenPort int    // Port to listen on (0 for random)
}
