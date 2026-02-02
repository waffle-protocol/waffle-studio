package wallet

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/waffle-studio/waffle/internal/config"
	"github.com/waffle-studio/waffle/internal/erc20"
)

// Wallet represents an Ethereum wallet connected to an RPC endpoint
type Wallet struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address
	client     *ethclient.Client
	rpcURL     string
}

// New creates a new wallet from private key and connects to RPC
func New(cfg *config.Config) (*Wallet, error) {
	if cfg.PrivateKey == "" {
		return nil, fmt.Errorf("private key not set")
	}
	// Remove 0x prefix if present
	keyHex := strings.TrimPrefix(cfg.PrivateKey, "0x")

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
	if cfg.RPCURL == "" {
		cfg.RPCURL = "http://127.0.0.1:8545"
	}
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	return &Wallet{
		privateKey: privateKey,
		address:    address,
		client:     client,
		rpcURL:     cfg.RPCURL,
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

// GetTokenBalance returns the ERC-20 token balance for the wallet
func (w *Wallet) GetTokenBalance(ctx context.Context, tokenAddress common.Address) (*big.Int, error) {
	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(erc20.BalanceOfABI))
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

// ToWei converts float amount to BigInt with given decimals
func ToWei(amount float64, decimals int) *big.Int {
	// Create multiplier 10^decimals
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	// Convert float to big.Float
	amountFl := big.NewFloat(amount)

	// Multiply
	resultFl := new(big.Float).Mul(amountFl, new(big.Float).SetInt(multiplier))

	// Convert to Int
	resultInt, _ := resultFl.Int(nil)
	return resultInt
}
