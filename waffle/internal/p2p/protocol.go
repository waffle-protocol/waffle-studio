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

// RequestHandler is called when a request is received from a peer
// It receives the prompt and file data, and returns the processed result
type RequestHandler func(prompt string, fileData []byte) ([]byte, error)

// SetupProvider sets up the node as a provider that handles incoming requests
func (n *Node) SetupProvider(handler RequestHandler) {
	n.Host.SetStreamHandler(ProtocolID, func(s network.Stream) {
		defer s.Close()

		// Read encrypted request
		reader := bufio.NewReader(s)
		encryptedData, err := io.ReadAll(reader)
		if err != nil {
			fmt.Printf("Error reading request: %v\n", err)
			return
		}

		// Decrypt request
		decrypted, err := Decrypt(encryptedData, n.encKey)
		if err != nil {
			fmt.Printf("Error decrypting request: %v\n", err)
			return
		}

		// Parse payload
		var payload SyrupPayload
		if err := json.Unmarshal(decrypted, &payload); err != nil {
			fmt.Printf("Error parsing payload: %v\n", err)
			return
		}

		// Process request
		result, err := handler(payload.Prompt, payload.Data)
		if err != nil {
			fmt.Printf("Error processing request: %v\n", err)
			// Send error response
			errPayload := SyrupPayload{Data: []byte(err.Error())}
			errData, _ := json.Marshal(errPayload)
			encrypted, _ := Encrypt(errData, n.encKey)
			s.Write(encrypted)
			return
		}

		// Prepare response payload
		responsePayload := SyrupPayload{Data: result}
		responseData, err := json.Marshal(responsePayload)
		if err != nil {
			fmt.Printf("Error marshaling response: %v\n", err)
			return
		}

		// Encrypt and send response
		encrypted, err := Encrypt(responseData, n.encKey)
		if err != nil {
			fmt.Printf("Error encrypting response: %v\n", err)
			return
		}

		s.Write(encrypted)
	})
}

// SendRequest sends a request to a peer and waits for the response
func (n *Node) SendRequest(ctx context.Context, peerID peer.ID, prompt string, fileData []byte) ([]byte, error) {
	// Connect to peer
	if err := n.Host.Connect(ctx, peer.AddrInfo{ID: peerID}); err != nil {
		return nil, fmt.Errorf("failed to connect to peer: %w", err)
	}

	// Open stream
	s, err := n.Host.NewStream(ctx, peerID, ProtocolID)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}
	defer s.Close()

	// Prepare request payload
	payload := SyrupPayload{
		Prompt: prompt,
		Data:   fileData,
	}
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Encrypt request
	encrypted, err := Encrypt(payloadData, n.encKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt request: %w", err)
	}

	// Send request and close write side
	if _, err := s.Write(encrypted); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	s.CloseWrite()

	// Read response
	reader := bufio.NewReader(s)
	encryptedResponse, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Decrypt response
	decrypted, err := Decrypt(encryptedResponse, n.encKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt response: %w", err)
	}

	// Parse response payload
	var responsePayload SyrupPayload
	if err := json.Unmarshal(decrypted, &responsePayload); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return responsePayload.Data, nil
}
