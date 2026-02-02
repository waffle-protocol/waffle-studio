package contracts

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/waffle-studio/waffle/internal/erc20"
)

// RegistryClient handles interaction with the BakeRegistry contract
type RegistryClient struct {
	address    common.Address
	client     *ethclient.Client
	privateKey []byte
	chainID    *big.Int
}

// RegistryABI is the ABI for BakeRegistry contract
const RegistryABI = `[
	{"inputs":[{"name":"codeHash","type":"bytes32"},{"name":"reward","type":"uint256"}],"name":"createRequest","outputs":[{"name":"requestId","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"name":"requestId","type":"uint256"}],"name":"getRequest","outputs":[{"components":[{"name":"requester","type":"address"},{"name":"baker","type":"address"},{"name":"codeHash","type":"bytes32"},{"name":"solutionHash","type":"bytes32"},{"name":"reward","type":"uint256"},{"name":"tokenUsage","type":"uint256"},{"name":"createdAt","type":"uint256"},{"name":"status","type":"uint8"}],"name":"","type":"tuple"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"nextRequestId","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
	{"inputs":[{"name":"requestId","type":"uint256"}],"name":"cancelRequest","outputs":[],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"name":"requestId","type":"uint256"},{"name":"paymentAmount","type":"uint256"}],"name":"acceptSolution","outputs":[],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"name":"requestId","type":"uint256"},{"name":"solutionHash","type":"bytes32"},{"name":"tokenUsage","type":"uint256"}],"name":"submitSolution","outputs":[],"stateMutability":"nonpayable","type":"function"}
]`


// NewRegistryClient creates a new BakeRegistry instance
func NewRegistryClient(address string, client *ethclient.Client, privateKeyHex string) (*RegistryClient, error) {
	keyHex := strings.TrimPrefix(privateKeyHex, "0x")
	privateKey, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	return &RegistryClient{
		address:    common.HexToAddress(address),
		client:     client,
		privateKey: crypto.FromECDSA(privateKey),
		chainID:    chainID,
	}, nil
}

// HashCode generates a bytes32 hash from file content
func HashCode(content string) [32]byte {
	hash := sha256.Sum256([]byte(content))
	return hash
}

// ApproveToken approves SYRUP tokens for the registry
func (r *RegistryClient) ApproveToken(ctx context.Context, tokenAddress string, amount *big.Int) (*types.Receipt, error) {
	parsedABI, err := abi.JSON(strings.NewReader(erc20.ApproveABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := parsedABI.Pack("approve", r.address, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to pack approve data: %w", err)
	}

	tokenAddr := common.HexToAddress(tokenAddress)
	return r.sendTransaction(ctx, &tokenAddr, data)
}

// CreateRequest creates a new bake request on-chain
func (r *RegistryClient) CreateRequest(ctx context.Context, codeHash [32]byte, reward *big.Int) (*types.Receipt, *big.Int, error) {
	parsedABI, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := parsedABI.Pack("createRequest", codeHash, reward)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pack createRequest data: %w", err)
	}

	receipt, err := r.sendTransaction(ctx, &r.address, data)
	if err != nil {
		return nil, nil, err
	}

	// Get request ID from nextRequestId - 1 (since it was incremented after creation)
	requestID, err := r.getNextRequestID(ctx)
	if err != nil {
		return receipt, nil, nil // Return receipt even if we can't get request ID
	}

	// Request ID is nextRequestId - 1
	actualID := new(big.Int).Sub(requestID, big.NewInt(1))
	return receipt, actualID, nil
}

// getNextRequestID fetches the next request ID from contract
func (r *RegistryClient) getNextRequestID(ctx context.Context) (*big.Int, error) {
	parsedABI, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return nil, err
	}

	data, err := parsedABI.Pack("nextRequestId")
	if err != nil {
		return nil, err
	}

	result, err := r.client.CallContract(ctx, ethereum.CallMsg{
		To:   &r.address,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}

	var nextID *big.Int
	err = parsedABI.UnpackIntoInterface(&nextID, "nextRequestId", result)
	if err != nil {
		return nil, err
	}

	return nextID, nil
}

// GetNextRequestID is the public version
func (r *RegistryClient) GetNextRequestID(ctx context.Context) (*big.Int, error) {
	return r.getNextRequestID(ctx)
}

// BakeRequest represents a request from the contract
type BakeRequest struct {
	Requester    common.Address
	Baker        common.Address
	CodeHash     [32]byte
	SolutionHash [32]byte
	Reward       *big.Int
	TokenUsage   *big.Int
	CreatedAt    *big.Int
	Status       uint8
}

// RequestStatus enum values
const (
	StatusPending   uint8 = 0
	StatusSubmitted uint8 = 1
	StatusAccepted  uint8 = 2
	StatusRejected  uint8 = 3
	StatusCancelled uint8 = 4
)

// GetRequest fetches a request by ID
func (r *RegistryClient) GetRequest(ctx context.Context, requestID *big.Int) (*BakeRequest, error) {
	parsedABI, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return nil, err
	}

	data, err := parsedABI.Pack("getRequest", requestID)
	if err != nil {
		return nil, err
	}

	result, err := r.client.CallContract(ctx, ethereum.CallMsg{
		To:   &r.address,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}

	// Unpack the tuple result using Unpack (returns []interface{})
	outputs, err := parsedABI.Unpack("getRequest", result)
	if err != nil {
		return nil, err
	}

	if len(outputs) == 0 {
		return nil, fmt.Errorf("no output from getRequest")
	}

	// The output is a struct, we need to parse it from the anonymous struct
	// outputs[0] is the tuple containing all fields
	tuple, ok := outputs[0].(struct {
		Requester    common.Address `json:"requester"`
		Baker        common.Address `json:"baker"`
		CodeHash     [32]byte       `json:"codeHash"`
		SolutionHash [32]byte       `json:"solutionHash"`
		Reward       *big.Int       `json:"reward"`
		TokenUsage   *big.Int       `json:"tokenUsage"`
		CreatedAt    *big.Int       `json:"createdAt"`
		Status       uint8          `json:"status"`
	})
	if !ok {
		return nil, fmt.Errorf("failed to parse request tuple")
	}

	return &BakeRequest{
		Requester:    tuple.Requester,
		Baker:        tuple.Baker,
		CodeHash:     tuple.CodeHash,
		SolutionHash: tuple.SolutionHash,
		Reward:       tuple.Reward,
		TokenUsage:   tuple.TokenUsage,
		CreatedAt:    tuple.CreatedAt,
		Status:       tuple.Status,
	}, nil
}

// CancelRequest cancels a pending request and refunds SYRUP
func (r *RegistryClient) CancelRequest(ctx context.Context, requestID *big.Int) (*types.Receipt, error) {
	parsedABI, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := parsedABI.Pack("cancelRequest", requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to pack cancelRequest data: %w", err)
	}

	return r.sendTransaction(ctx, &r.address, data)
}

// AcceptSolution accepts a submitted solution and pays the baker
func (r *RegistryClient) AcceptSolution(ctx context.Context, requestID *big.Int, paymentAmount *big.Int) (*types.Receipt, error) {
	parsedABI, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := parsedABI.Pack("acceptSolution", requestID, paymentAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to pack acceptSolution data: %w", err)
	}

	return r.sendTransaction(ctx, &r.address, data)
}

// SubmitSolution submits a solution for a pending request (Baker function)
func (r *RegistryClient) SubmitSolution(ctx context.Context, requestID *big.Int, solutionHash [32]byte, tokenUsage *big.Int) (*types.Receipt, error) {
	parsedABI, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := parsedABI.Pack("submitSolution", requestID, solutionHash, tokenUsage)
	if err != nil {
		return nil, fmt.Errorf("failed to pack submitSolution data: %w", err)
	}

	return r.sendTransaction(ctx, &r.address, data)
}

// sendTransaction sends a signed transaction
func (r *RegistryClient) sendTransaction(ctx context.Context, to *common.Address, data []byte) (*types.Receipt, error) {
	privateKey, err := crypto.ToECDSA(r.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := r.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := r.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	// Estimate gas
	gasLimit, err := r.client.EstimateGas(ctx, ethereum.CallMsg{
		From: fromAddress,
		To:   to,
		Data: data,
	})
	if err != nil {
		gasLimit = 500000 // Default gas limit if estimation fails
	}

	tx := types.NewTransaction(nonce, *to, big.NewInt(0), gasLimit, gasPrice, data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(r.chainID), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = r.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	// Wait for receipt
	receipt, err := waitForReceipt(ctx, r.client, signedTx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get receipt: %w", err)
	}

	return receipt, nil
}

// waitForReceipt waits for transaction to be mined
func waitForReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Keep polling
		}
	}
}
