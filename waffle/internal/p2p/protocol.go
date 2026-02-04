package p2p

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

// ProcessResult contains the result of processing a request
type ProcessResult struct {
	Data       []byte // Processed data
	TokenUsage uint64 // Number of AI tokens used
}

// RequestHandler is called when a request is received from a peer
// It receives the prompt and file data, and returns the processed result with token usage
type RequestHandler func(prompt string, fileData []byte) (*ProcessResult, error)

// SetupProvider sets up the node as a provider that handles incoming requests
// providerAddress is the provider's wallet address for blockchain payments
func (n *Node) SetupProvider(providerAddress string, handler RequestHandler) {
	n.Host.SetStreamHandler(ProtocolID, func(s network.Stream) {
		defer s.Close()

		// Read request (libp2p already encrypts the stream via TLS/Noise)
		reader := bufio.NewReader(s)
		data, err := io.ReadAll(reader)
		if err != nil {
			fmt.Printf("Error reading request: %v\n", err)
			return
		}

		// Parse payload
		var payload SyrupPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			fmt.Printf("Error parsing payload: %v\n", err)
			return
		}

		// Process request
		result, err := handler(payload.Prompt, payload.Data)
		if err != nil {
			fmt.Printf("Error processing request: %v\n", err)
			// Send error response
			errPayload := SyrupPayload{
				Error:           err.Error(),
				ProviderAddress: providerAddress,
			}
			errData, _ := json.Marshal(errPayload)
			s.Write(errData)
			return
		}

		// Prepare and send response with provider address and token usage
		responsePayload := SyrupPayload{
			Data:            result.Data,
			ProviderAddress: providerAddress,
			TokenUsage:      result.TokenUsage,
		}
		responseData, err := json.Marshal(responsePayload)
		if err != nil {
			fmt.Printf("Error marshaling response: %v\n", err)
			return
		}

		s.Write(responseData)
	})
}

// SendRequest sends a request to a peer and waits for the response
// Returns the full SyrupPayload containing data, provider address, and token usage
func (n *Node) SendRequest(ctx context.Context, peerInfo peer.AddrInfo, prompt string, fileData []byte) (*SyrupPayload, error) {
	return n.SendRequestWithID(ctx, peerInfo, prompt, fileData, 0)
}

// SendRequestWithID sends a request with a blockchain request ID to a peer
func (n *Node) SendRequestWithID(ctx context.Context, peerInfo peer.AddrInfo, prompt string, fileData []byte, requestID uint64) (*SyrupPayload, error) {
	// Connect to peer using full address info
	if err := n.Host.Connect(ctx, peerInfo); err != nil {
		return nil, fmt.Errorf("failed to connect to peer: %w", err)
	}

	// Open stream (libp2p already encrypts via TLS/Noise)
	s, err := n.Host.NewStream(ctx, peerInfo.ID, ProtocolID)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}
	defer s.Close()

	// Prepare request payload
	payload := SyrupPayload{
		Prompt:    prompt,
		Data:      fileData,
		RequestID: requestID,
	}
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Send request and close write side
	if _, err := s.Write(payloadData); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	s.CloseWrite()

	// Read response
	reader := bufio.NewReader(s)
	responseData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response payload
	var responsePayload SyrupPayload
	if err := json.Unmarshal(responseData, &responsePayload); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for error in response
	if responsePayload.Error != "" {
		return nil, fmt.Errorf("provider error: %s", responsePayload.Error)
	}

	return &responsePayload, nil
}
