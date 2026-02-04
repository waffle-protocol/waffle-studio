package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Config holds relay server configuration
type Config struct {
	RPCURL       string
	PrivateKey   string
	BakeRegistry string
	SyrupToken   string
	Port         string
}

var (
	client       *ethclient.Client
	privateKey   []byte
	chainID      *big.Int
	config       Config
	relayAddress common.Address
)

// ABI definitions
const RegistryABI = `[
	{"inputs":[{"name":"codeHash","type":"bytes32"},{"name":"reward","type":"uint256"}],"name":"createRequest","outputs":[{"name":"requestId","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"name":"requestId","type":"uint256"}],"name":"cancelRequest","outputs":[],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"name":"requestId","type":"uint256"},{"name":"paymentAmount","type":"uint256"}],"name":"acceptSolution","outputs":[],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"name":"requestId","type":"uint256"}],"name":"rejectSolution","outputs":[],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"name":"requestId","type":"uint256"},{"name":"solutionHash","type":"bytes32"},{"name":"tokenUsage","type":"uint256"}],"name":"submitSolution","outputs":[],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[],"name":"nextRequestId","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
	{"inputs":[{"name":"codeHash","type":"bytes32"},{"name":"reward","type":"uint256"},{"name":"requester","type":"address"},{"name":"signature","type":"bytes"}],"name":"createRequestGasless","outputs":[{"name":"requestId","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"name":"requestId","type":"uint256"},{"name":"requester","type":"address"},{"name":"signature","type":"bytes"}],"name":"cancelRequestGasless","outputs":[],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"name":"requestId","type":"uint256"},{"name":"paymentAmount","type":"uint256"},{"name":"requester","type":"address"},{"name":"signature","type":"bytes"}],"name":"acceptSolutionGasless","outputs":[],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"name":"requestId","type":"uint256"},{"name":"requester","type":"address"},{"name":"signature","type":"bytes"}],"name":"rejectSolutionGasless","outputs":[],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"name":"requestId","type":"uint256"},{"name":"solutionHash","type":"bytes32"},{"name":"tokenUsage","type":"uint256"},{"name":"baker","type":"address"},{"name":"signature","type":"bytes"}],"name":"submitSolutionGasless","outputs":[],"stateMutability":"nonpayable","type":"function"}
]`

const ERC20ABI = `[
	{"inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}
]`

// Request/Response types

type FaucetRequest struct {
	Address string `json:"address"`
}

type FaucetResponse struct {
	Success     bool   `json:"success"`
	TxHash      string `json:"tx_hash,omitempty"`
	EthAmount   string `json:"eth_amount,omitempty"`
	SyrupAmount string `json:"syrup_amount,omitempty"`
	Error       string `json:"error,omitempty"`
}

type CreateRequestReq struct {
	Requester string `json:"requester"`
	CodeHash  string `json:"code_hash"`
	Reward    string `json:"reward"`
	Signature string `json:"signature"`
}

type CreateRequestResp struct {
	Success   bool   `json:"success"`
	RequestID uint64 `json:"request_id,omitempty"`
	TxHash    string `json:"tx_hash,omitempty"`
	Error     string `json:"error,omitempty"`
}

type AcceptSolutionReq struct {
	Requester     string `json:"requester"`
	RequestID     uint64 `json:"request_id"`
	PaymentAmount string `json:"payment_amount"`
	Signature     string `json:"signature"`
}

type RejectSolutionReq struct {
	Requester string `json:"requester"`
	RequestID uint64 `json:"request_id"`
	Signature string `json:"signature"`
}

type CancelRequestReq struct {
	Requester string `json:"requester"`
	RequestID uint64 `json:"request_id"`
	Signature string `json:"signature"`
}

type SubmitSolutionReq struct {
	RequestID    uint64 `json:"request_id"`
	SolutionHash string `json:"solution_hash"`
	TokenUsage   uint64 `json:"token_usage"`
	Signature    string `json:"signature"`
	Signer       string `json:"signer"`
}

type GenericResponse struct {
	Success bool   `json:"success"`
	TxHash  string `json:"tx_hash,omitempty"`
	Error   string `json:"error,omitempty"`
}

func main() {
	// Load configuration from environment
	config = Config{
		RPCURL:       getEnv("RPC_URL", "https://sepolia.base.org"),
		PrivateKey:   getEnv("RELAY_PRIVATE_KEY", ""),
		BakeRegistry: getEnv("BAKE_REGISTRY", "0x8c03552b7ae490ddc2e6b5a1b3452129e1135323"),
		SyrupToken:   getEnv("SYRUP_TOKEN", "0xb84284ddab9f7e2b14ba81c6f44db8d99488e23b"),
		Port:         getEnv("PORT", "8080"),
	}

	if config.PrivateKey == "" {
		log.Fatal("RELAY_PRIVATE_KEY environment variable is required")
	}

	// Initialize Ethereum client
	var err error
	client, err = ethclient.Dial(config.RPCURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum: %v", err)
	}

	// Parse private key
	keyHex := strings.TrimPrefix(config.PrivateKey, "0x")
	pk, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}
	privateKey = crypto.FromECDSA(pk)
	relayAddress = crypto.PubkeyToAddress(pk.PublicKey)

	// Get chain ID
	chainID, err = client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Setup HTTP routes with CORS
	mux := http.NewServeMux()
	mux.HandleFunc("/health", corsMiddleware(healthHandler))
	mux.HandleFunc("/faucet", corsMiddleware(faucetHandler))
	mux.HandleFunc("/create-request", corsMiddleware(createRequestHandler))
	mux.HandleFunc("/accept-solution", corsMiddleware(acceptSolutionHandler))
	mux.HandleFunc("/reject-solution", corsMiddleware(rejectSolutionHandler))
	mux.HandleFunc("/cancel-request", corsMiddleware(cancelRequestHandler))
	mux.HandleFunc("/submit-solution", corsMiddleware(submitSolutionHandler))

	fmt.Println()
	fmt.Println("🧇 Waffle Gasless Relay Server")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("   Network:  Base Sepolia\n")
	fmt.Printf("   RPC:      %s\n", config.RPCURL)
	fmt.Printf("   Relay:    %s\n", relayAddress.Hex())
	fmt.Printf("   Registry: %s\n", config.BakeRegistry)
	fmt.Printf("   SYRUP:    %s\n", config.SyrupToken)
	fmt.Printf("   Port:     %s\n", config.Port)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("Endpoints:")
	fmt.Println("  POST /faucet          - Get free ETH + SYRUP")
	fmt.Println("  POST /create-request  - Create bake request (gasless)")
	fmt.Println("  POST /accept-solution - Accept solution (gasless)")
	fmt.Println("  POST /reject-solution - Reject solution (gasless)")
	fmt.Println("  POST /cancel-request  - Cancel request (gasless)")
	fmt.Println("  POST /submit-solution - Submit solution (gasless)")
	fmt.Println()

	log.Printf("Relay server listening on :%s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, mux))
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "relay": relayAddress.Hex()})
}

// ============================================================================
// Faucet Handler - Sends ETH + SYRUP to user
// ============================================================================

func faucetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FaucetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !common.IsHexAddress(req.Address) {
		sendError(w, "Invalid address", http.StatusBadRequest)
		return
	}

	userAddr := common.HexToAddress(req.Address)
	log.Printf("[Faucet] Sending to %s", userAddr.Hex())

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Send ETH (0.001 ETH for gas)
	ethAmount := big.NewInt(1e15) // 0.001 ETH
	ethTxHash, err := sendETH(ctx, userAddr, ethAmount)
	if err != nil {
		log.Printf("[Faucet] ETH transfer failed: %v", err)
		sendError(w, fmt.Sprintf("ETH transfer failed: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("[Faucet] ETH sent: %s", ethTxHash)

	// Send SYRUP (100 SYRUP)
	syrupAmount := new(big.Int)
	syrupAmount.SetString("100000000000000000000", 10) // 100 * 10^18
	syrupTxHash, err := sendSYRUP(ctx, userAddr, syrupAmount)
	if err != nil {
		log.Printf("[Faucet] SYRUP transfer failed: %v", err)
		// Still return success for ETH
		json.NewEncoder(w).Encode(FaucetResponse{
			Success:     true,
			TxHash:      ethTxHash,
			EthAmount:   "0.001",
			SyrupAmount: "0 (transfer failed)",
			Error:       fmt.Sprintf("SYRUP transfer failed: %v", err),
		})
		return
	}
	log.Printf("[Faucet] SYRUP sent: %s", syrupTxHash)

	json.NewEncoder(w).Encode(FaucetResponse{
		Success:     true,
		TxHash:      syrupTxHash,
		EthAmount:   "0.001",
		SyrupAmount: "100",
	})
}

func sendETH(ctx context.Context, to common.Address, amount *big.Int) (string, error) {
	pk, _ := crypto.ToECDSA(privateKey)

	nonce, err := client.PendingNonceAt(ctx, relayAddress)
	if err != nil {
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(nonce, to, amount, 21000, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), pk)
	if err != nil {
		return "", err
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

func sendSYRUP(ctx context.Context, to common.Address, amount *big.Int) (string, error) {
	parsedABI, _ := abi.JSON(strings.NewReader(ERC20ABI))
	data, err := parsedABI.Pack("transfer", to, amount)
	if err != nil {
		return "", err
	}

	tokenAddr := common.HexToAddress(config.SyrupToken)
	return sendTx(ctx, &tokenAddr, data, nil)
}

// ============================================================================
// Create Request Handler (Gasless)
// ============================================================================

func createRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateRequestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[CreateRequest] From %s, reward %s", req.Requester, req.Reward)

	// Parse code hash
	codeHash, err := parseBytes32(req.CodeHash)
	if err != nil {
		sendError(w, "Invalid code hash", http.StatusBadRequest)
		return
	}

	// Parse reward
	reward := new(big.Int)
	reward.SetString(req.Reward, 10)

	// Parse signature
	signature, err := hex.DecodeString(strings.TrimPrefix(req.Signature, "0x"))
	if err != nil {
		sendError(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	requesterAddr := common.HexToAddress(req.Requester)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Get next request ID before creating
	nextID, err := getNextRequestID(ctx)
	if err != nil {
		log.Printf("[CreateRequest] Failed to get next ID: %v", err)
	}

	// Call createRequestGasless
	parsedABI, _ := abi.JSON(strings.NewReader(RegistryABI))
	data, err := parsedABI.Pack("createRequestGasless", codeHash, reward, requesterAddr, signature)
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to pack data: %v", err), http.StatusInternalServerError)
		return
	}

	registryAddr := common.HexToAddress(config.BakeRegistry)
	txHash, err := sendTx(ctx, &registryAddr, data, nil)
	if err != nil {
		log.Printf("[CreateRequest] Failed: %v", err)
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[CreateRequest] Success: txHash=%s, requestID=%d", txHash, nextID)

	json.NewEncoder(w).Encode(CreateRequestResp{
		Success:   true,
		RequestID: nextID,
		TxHash:    txHash,
	})
}

func getNextRequestID(ctx context.Context) (uint64, error) {
	parsedABI, _ := abi.JSON(strings.NewReader(RegistryABI))
	data, _ := parsedABI.Pack("nextRequestId")

	registryAddr := common.HexToAddress(config.BakeRegistry)
	result, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &registryAddr,
		Data: data,
	}, nil)
	if err != nil {
		return 0, err
	}

	var nextID *big.Int
	parsedABI.UnpackIntoInterface(&nextID, "nextRequestId", result)
	return nextID.Uint64(), nil
}

// ============================================================================
// Accept Solution Handler (Gasless)
// ============================================================================

func acceptSolutionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AcceptSolutionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[AcceptSolution] RequestID=%d, Requester=%s", req.RequestID, req.Requester)

	paymentAmount := new(big.Int)
	paymentAmount.SetString(req.PaymentAmount, 10)

	signature, err := hex.DecodeString(strings.TrimPrefix(req.Signature, "0x"))
	if err != nil {
		sendError(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	requesterAddr := common.HexToAddress(req.Requester)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	parsedABI, _ := abi.JSON(strings.NewReader(RegistryABI))
	data, err := parsedABI.Pack("acceptSolutionGasless", big.NewInt(int64(req.RequestID)), paymentAmount, requesterAddr, signature)
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to pack data: %v", err), http.StatusInternalServerError)
		return
	}

	registryAddr := common.HexToAddress(config.BakeRegistry)
	txHash, err := sendTx(ctx, &registryAddr, data, nil)
	if err != nil {
		log.Printf("[AcceptSolution] Failed: %v", err)
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[AcceptSolution] Success: %s", txHash)

	json.NewEncoder(w).Encode(GenericResponse{
		Success: true,
		TxHash:  txHash,
	})
}

// ============================================================================
// Reject Solution Handler (Gasless)
// ============================================================================

func rejectSolutionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RejectSolutionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[RejectSolution] RequestID=%d, Requester=%s", req.RequestID, req.Requester)

	signature, err := hex.DecodeString(strings.TrimPrefix(req.Signature, "0x"))
	if err != nil {
		sendError(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	requesterAddr := common.HexToAddress(req.Requester)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	parsedABI, _ := abi.JSON(strings.NewReader(RegistryABI))
	data, err := parsedABI.Pack("rejectSolutionGasless", big.NewInt(int64(req.RequestID)), requesterAddr, signature)
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to pack data: %v", err), http.StatusInternalServerError)
		return
	}

	registryAddr := common.HexToAddress(config.BakeRegistry)
	txHash, err := sendTx(ctx, &registryAddr, data, nil)
	if err != nil {
		log.Printf("[RejectSolution] Failed: %v", err)
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[RejectSolution] Success: %s", txHash)

	json.NewEncoder(w).Encode(GenericResponse{
		Success: true,
		TxHash:  txHash,
	})
}

// ============================================================================
// Cancel Request Handler (Gasless)
// ============================================================================

func cancelRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CancelRequestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[CancelRequest] RequestID=%d, Requester=%s", req.RequestID, req.Requester)

	signature, err := hex.DecodeString(strings.TrimPrefix(req.Signature, "0x"))
	if err != nil {
		sendError(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	requesterAddr := common.HexToAddress(req.Requester)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	parsedABI, _ := abi.JSON(strings.NewReader(RegistryABI))
	data, err := parsedABI.Pack("cancelRequestGasless", big.NewInt(int64(req.RequestID)), requesterAddr, signature)
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to pack data: %v", err), http.StatusInternalServerError)
		return
	}

	registryAddr := common.HexToAddress(config.BakeRegistry)
	txHash, err := sendTx(ctx, &registryAddr, data, nil)
	if err != nil {
		log.Printf("[CancelRequest] Failed: %v", err)
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[CancelRequest] Success: %s", txHash)

	json.NewEncoder(w).Encode(GenericResponse{
		Success: true,
		TxHash:  txHash,
	})
}

// ============================================================================
// Submit Solution Handler (Gasless) - For providers
// ============================================================================

func submitSolutionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SubmitSolutionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[SubmitSolution] RequestID=%d, Signer=%s", req.RequestID, req.Signer)

	solutionHash, err := parseBytes32(req.SolutionHash)
	if err != nil {
		sendError(w, "Invalid solution hash", http.StatusBadRequest)
		return
	}

	signature, err := hex.DecodeString(strings.TrimPrefix(req.Signature, "0x"))
	if err != nil {
		sendError(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	signerAddr := common.HexToAddress(req.Signer)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	parsedABI, _ := abi.JSON(strings.NewReader(RegistryABI))
	data, err := parsedABI.Pack("submitSolutionGasless", big.NewInt(int64(req.RequestID)), solutionHash, big.NewInt(int64(req.TokenUsage)), signerAddr, signature)
	if err != nil {
		sendError(w, fmt.Sprintf("Failed to pack data: %v", err), http.StatusInternalServerError)
		return
	}

	registryAddr := common.HexToAddress(config.BakeRegistry)
	txHash, err := sendTx(ctx, &registryAddr, data, nil)
	if err != nil {
		log.Printf("[SubmitSolution] Failed: %v", err)
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[SubmitSolution] Success: %s", txHash)

	json.NewEncoder(w).Encode(GenericResponse{
		Success: true,
		TxHash:  txHash,
	})
}

// ============================================================================
// Helper Functions
// ============================================================================

func sendTx(ctx context.Context, to *common.Address, data []byte, value *big.Int) (string, error) {
	pk, _ := crypto.ToECDSA(privateKey)

	nonce, err := client.PendingNonceAt(ctx, relayAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	if value == nil {
		value = big.NewInt(0)
	}

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  relayAddress,
		To:    to,
		Data:  data,
		Value: value,
	})
	if err != nil {
		gasLimit = 500000 // Default
	}

	tx := types.NewTransaction(nonce, *to, value, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), pk)
	if err != nil {
		return "", fmt.Errorf("failed to sign tx: %w", err)
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return "", fmt.Errorf("failed to send tx: %w", err)
	}

	// Wait for receipt
	for {
		receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			if receipt.Status == 0 {
				return "", fmt.Errorf("transaction reverted")
			}
			return signedTx.Hash().Hex(), nil
		}

		select {
		case <-ctx.Done():
			return signedTx.Hash().Hex(), nil
		case <-time.After(500 * time.Millisecond):
			continue
		}
	}
}

func parseBytes32(s string) ([32]byte, error) {
	var result [32]byte
	hashStr := strings.TrimPrefix(s, "0x")
	if len(hashStr) != 64 {
		return result, fmt.Errorf("invalid length")
	}
	hashBytes, err := hex.DecodeString(hashStr)
	if err != nil {
		return result, err
	}
	copy(result[:], hashBytes)
	return result, nil
}

func sendError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(GenericResponse{
		Success: false,
		Error:   message,
	})
}
