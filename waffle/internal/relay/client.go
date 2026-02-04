package relay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client handles communication with the gasless relay server
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new relay client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// FaucetRequest represents a faucet claim request
type FaucetRequest struct {
	Address string `json:"address"`
}

// FaucetResponse represents the faucet response
type FaucetResponse struct {
	Success    bool   `json:"success"`
	TxHash     string `json:"tx_hash,omitempty"`
	EthAmount  string `json:"eth_amount,omitempty"`
	SyrupAmount string `json:"syrup_amount,omitempty"`
	Error      string `json:"error,omitempty"`
}

// CreateRequestReq represents a create request payload
type CreateRequestReq struct {
	Requester string `json:"requester"`
	CodeHash  string `json:"code_hash"`
	Reward    string `json:"reward"` // Wei as string
	Signature string `json:"signature"`
}

// CreateRequestResp represents the create request response
type CreateRequestResp struct {
	Success   bool   `json:"success"`
	RequestID uint64 `json:"request_id,omitempty"`
	TxHash    string `json:"tx_hash,omitempty"`
	Error     string `json:"error,omitempty"`
}

// AcceptSolutionReq represents an accept solution request
type AcceptSolutionReq struct {
	Requester     string `json:"requester"`
	RequestID     uint64 `json:"request_id"`
	PaymentAmount string `json:"payment_amount"` // Wei as string
	Signature     string `json:"signature"`
}

// RejectSolutionReq represents a reject solution request
type RejectSolutionReq struct {
	Requester string `json:"requester"`
	RequestID uint64 `json:"request_id"`
	Signature string `json:"signature"`
}

// CancelRequestReq represents a cancel request payload
type CancelRequestReq struct {
	Requester string `json:"requester"`
	RequestID uint64 `json:"request_id"`
	Signature string `json:"signature"`
}

// GenericResponse represents a generic relay response
type GenericResponse struct {
	Success bool   `json:"success"`
	TxHash  string `json:"tx_hash,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Faucet claims ETH and SYRUP tokens from the relay faucet
func (c *Client) Faucet(address string) (*FaucetResponse, error) {
	req := FaucetRequest{Address: address}

	resp, err := c.post("/faucet", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result FaucetResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("faucet failed: %s", result.Error)
	}

	return &result, nil
}

// CreateRequest creates a new bake request via relay
func (c *Client) CreateRequest(requester, codeHash, reward, signature string) (*CreateRequestResp, error) {
	req := CreateRequestReq{
		Requester: requester,
		CodeHash:  codeHash,
		Reward:    reward,
		Signature: signature,
	}

	resp, err := c.post("/create-request", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CreateRequestResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("create request failed: %s", result.Error)
	}

	return &result, nil
}

// AcceptSolution accepts a solution via relay
func (c *Client) AcceptSolution(requester string, requestID uint64, paymentAmount, signature string) (*GenericResponse, error) {
	req := AcceptSolutionReq{
		Requester:     requester,
		RequestID:     requestID,
		PaymentAmount: paymentAmount,
		Signature:     signature,
	}

	resp, err := c.post("/accept-solution", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GenericResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("accept solution failed: %s", result.Error)
	}

	return &result, nil
}

// RejectSolution rejects a solution via relay
func (c *Client) RejectSolution(requester string, requestID uint64, signature string) (*GenericResponse, error) {
	req := RejectSolutionReq{
		Requester: requester,
		RequestID: requestID,
		Signature: signature,
	}

	resp, err := c.post("/reject-solution", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GenericResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("reject solution failed: %s", result.Error)
	}

	return &result, nil
}

// CancelRequest cancels a pending request via relay
func (c *Client) CancelRequest(requester string, requestID uint64, signature string) (*GenericResponse, error) {
	req := CancelRequestReq{
		Requester: requester,
		RequestID: requestID,
		Signature: signature,
	}

	resp, err := c.post("/cancel-request", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GenericResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("cancel request failed: %s", result.Error)
	}

	return &result, nil
}

// Health checks if the relay server is healthy
func (c *Client) Health() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("relay server unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("relay server unhealthy: status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) post(path string, body interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+path, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
