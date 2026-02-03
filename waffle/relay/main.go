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

// SubmitRequest represents a gasless submit request
type SubmitRequest struct {
	RequestID    uint64 `json:"request_id"`
	SolutionHash string `json:"solution_hash"`
	TokenUsage   uint64 `json:"token_usage"`
	Signature    string `json:"signature"`
	Signer       string `json:"signer"`
}

// RelayResponse represents the response from relay
type RelayResponse struct {
	Success bool   `json:"success"`
	TxHash  string `json:"tx_hash,omitempty"`
	Error   string `json:"error,omitempty"`
}

var (
	client     *ethclient.Client
	privateKey []byte
	chainID    *big.Int
	config     Config
)

const RegistryABI = `[
	{"inputs":[{"name":"requestId","type":"uint256"},{"name":"solutionHash","type":"bytes32"},{"name":"tokenUsage","type":"uint256"}],"name":"submitSolution","outputs":[],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"name":"requestId","type":"uint256"},{"name":"solutionHash","type":"bytes32"},{"name":"tokenUsage","type":"uint256"},{"name":"baker","type":"address"},{"name":"signature","type":"bytes"}],"name":"submitSolutionGasless","outputs":[],"stateMutability":"nonpayable","type":"function"}
]`

func main() {
	// Load configuration from environment
	config = Config{
		RPCURL:       getEnv("RPC_URL", "https://sepolia.base.org"),
		PrivateKey:   getEnv("RELAY_PRIVATE_KEY", ""),
		BakeRegistry: getEnv("BAKE_REGISTRY", "0xf53ce337d2ec220e5c166be826bf14f8ef6da36e"),
		SyrupToken:   getEnv("SYRUP_TOKEN", "0x497a760b88551f9d1957493d000651f0d473360b"),
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

	// Get chain ID
	chainID, err = client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Setup HTTP routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/submit-solution", submitSolutionHandler)

	// Get relay address for logging
	relayAddress := crypto.PubkeyToAddress(pk.PublicKey)

	fmt.Println()
	fmt.Println("🧇 Waffle Relay Server")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("   Network:  %s\n", config.RPCURL)
	fmt.Printf("   Relay:    %s\n", relayAddress.Hex())
	fmt.Printf("   Port:     %s\n", config.Port)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	log.Printf("Relay server listening on :%s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func submitSolutionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received submit request: requestID=%d, signer=%s", req.RequestID, req.Signer)

	// Parse solution hash
	hashStr := strings.TrimPrefix(req.SolutionHash, "0x")
	if len(hashStr) != 64 {
		sendError(w, "Invalid solution hash length", http.StatusBadRequest)
		return
	}

	hashBytes, err := hex.DecodeString(hashStr)
	if err != nil {
		sendError(w, "Invalid solution hash format", http.StatusBadRequest)
		return
	}

	var solutionHash [32]byte
	copy(solutionHash[:], hashBytes)

	// Parse signature
	signatureBytes, err := hex.DecodeString(req.Signature)
	if err != nil {
		sendError(w, "Invalid signature format", http.StatusBadRequest)
		return
	}

	// Submit transaction
	signerAddr := common.HexToAddress(req.Signer)
	txHash, err := submitSolutionTx(req.RequestID, solutionHash, req.TokenUsage, signerAddr, signatureBytes)
	if err != nil {
		log.Printf("Submit failed: %v", err)
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Solution submitted: txHash=%s", txHash)

	json.NewEncoder(w).Encode(RelayResponse{
		Success: true,
		TxHash:  txHash,
	})
}

func submitSolutionTx(requestID uint64, solutionHash [32]byte, tokenUsage uint64, signer common.Address, signature []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	parsedABI, err := abi.JSON(strings.NewReader(RegistryABI))
	if err != nil {
		return "", fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := parsedABI.Pack("submitSolutionGasless", big.NewInt(int64(requestID)), solutionHash, big.NewInt(int64(tokenUsage)), signer, signature)
	if err != nil {
		return "", fmt.Errorf("failed to pack data: %w", err)
	}

	// ... (rest is same)

	pk, err := crypto.ToECDSA(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	fromAddress := crypto.PubkeyToAddress(pk.PublicKey)
	registryAddress := common.HexToAddress(config.BakeRegistry)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From: fromAddress,
		To:   &registryAddress,
		Data: data,
	})
	if err != nil {
		gasLimit = 200000
	}

	tx := types.NewTransaction(nonce, registryAddress, big.NewInt(0), gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), pk)
	if err != nil {
		return "", fmt.Errorf("failed to sign tx: %w", err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
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
			return signedTx.Hash().Hex(), nil // Return hash even if not confirmed yet
		case <-time.After(500 * time.Millisecond):
			continue
		}
	}
}

func sendError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(RelayResponse{
		Success: false,
		Error:   message,
	})
}
