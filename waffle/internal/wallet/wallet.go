package wallet

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Wallet represents an Ethereum wallet connected to an RPC endpoint
type Wallet struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address
	client     *ethclient.Client
	rpcURL     string
}

// Config holds wallet configuration from environment
type Config struct {
	PrivateKey string
	RPCURL     string
	SyrupToken string
}

// LoadConfig loads wallet configuration from environment variables
func LoadConfig() (*Config, error) {
	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		return nil, fmt.Errorf("PRIVATE_KEY environment variable not set")
	}

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		rpcURL = "http://127.0.0.1:8545" // Default to local Anvil
	}

	syrupToken := os.Getenv("SYRUP_TOKEN")
	if syrupToken == "" {
		return nil, fmt.Errorf("SYRUP_TOKEN environment variable not set")
	}

	return &Config{
		PrivateKey: privateKey,
		RPCURL:     rpcURL,
		SyrupToken: syrupToken,
	}, nil
}

// New creates a new wallet from private key and connects to RPC
func New(config *Config) (*Wallet, error) {
	// Remove 0x prefix if present
	keyHex := strings.TrimPrefix(config.PrivateKey, "0x")

	// Parse private key
	privateKey, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// Derive address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to get public key")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Connect to RPC
	client, err := ethclient.Dial(config.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	return &Wallet{
		privateKey: privateKey,
		address:    address,
		client:     client,
		rpcURL:     config.RPCURL,
	}, nil
}

// Address returns the wallet's Ethereum address
func (w *Wallet) Address() common.Address {
	return w.address
}

// AddressHex returns the wallet's address as hex string
func (w *Wallet) AddressHex() string {
	return w.address.Hex()
}

// Client returns the underlying ethclient
func (w *Wallet) Client() *ethclient.Client {
	return w.client
}

// Close closes the wallet's RPC connection
func (w *Wallet) Close() {
	if w.client != nil {
		w.client.Close()
	}
}

// GetETHBalance returns the ETH balance of the wallet
func (w *Wallet) GetETHBalance(ctx context.Context) (*big.Int, error) {
	return w.client.BalanceAt(ctx, w.address, nil)
}

// ERC20 ABI for balanceOf
const erc20BalanceOfABI = `[{"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`

// GetTokenBalance returns the ERC-20 token balance for the wallet
func (w *Wallet) GetTokenBalance(ctx context.Context, tokenAddress common.Address) (*big.Int, error) {
	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(erc20BalanceOfABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Encode function call
	data, err := parsedABI.Pack("balanceOf", w.address)
	if err != nil {
		return nil, fmt.Errorf("failed to pack call data: %w", err)
	}

	// Make call
	result, err := w.client.CallContract(ctx, ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	// Decode result
	var balance *big.Int
	err = parsedABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %w", err)
	}

	return balance, nil
}

// FormatTokenBalance converts raw balance to human-readable format
func FormatTokenBalance(balance *big.Int, decimals int) float64 {
	if balance == nil {
		return 0
	}

	// Create divisor (10^decimals)
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	// Convert to float
	balanceFloat := new(big.Float).SetInt(balance)
	divisorFloat := new(big.Float).SetInt(divisor)

	result := new(big.Float).Quo(balanceFloat, divisorFloat)
	f, _ := result.Float64()

	return f
}
