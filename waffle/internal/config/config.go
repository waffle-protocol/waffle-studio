package config

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/yaml.v3"
)

// Default testnet configuration
const (
	DefaultRPCURL       = "https://sepolia.base.org"
	DefaultSyrupToken   = "0xb84284ddab9f7e2b14ba81c6f44db8d99488e23b"
	DefaultBakeRegistry = "0x8c03552b7ae490ddc2e6b5a1b3452129e1135323"
	DefaultRelayURL     = "https://relay.waffle.studio" // Gasless relay server
)

type Config struct {
	RPCURL       string `yaml:"rpc_url"`
	PrivateKey   string `yaml:"private_key"`
	SyrupToken   string `yaml:"syrup_token"`
	BakeRegistry string `yaml:"bake_registry"`
	RelayURL     string `yaml:"relay_url"`
	PricingRate  uint64 `yaml:"pricing_rate"` // Rate per character+usage
}

// GeneratePrivateKey creates a new random Ethereum private key
func GeneratePrivateKey() (string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	// Convert to hex string
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	// Get address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", fmt.Errorf("failed to get public key")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return privateKeyHex, address, nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".waffle", "config.yaml"), nil
}

// EnsureConfig checks if config exists, if not creates one with auto-generated wallet
func EnsureConfig() (*Config, bool, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, false, err
	}

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		// Config exists, load it
		cfg, err := loadFromFile(configPath)
		return cfg, false, err
	}

	// Config doesn't exist, create new one with generated wallet
	privateKey, address, err := GeneratePrivateKey()
	if err != nil {
		return nil, false, err
	}

	cfg := &Config{
		RPCURL:       DefaultRPCURL,
		PrivateKey:   privateKey,
		SyrupToken:   DefaultSyrupToken,
		BakeRegistry: DefaultBakeRegistry,
		RelayURL:     DefaultRelayURL,
		PricingRate:  1,
	}

	// Create directory
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, false, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write config file
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, false, fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add header comment
	header := fmt.Sprintf(`# Waffle CLI Configuration
# Auto-generated wallet address: %s
# Network: Base Sepolia Testnet
#
# ⚠️  IMPORTANT: Keep your private_key secret!
# 💡 Get testnet ETH from: https://www.alchemy.com/faucets/base-sepolia

`, address)

	if err := os.WriteFile(configPath, []byte(header+string(data)), 0600); err != nil {
		return nil, false, fmt.Errorf("failed to write config file: %w", err)
	}

	// Print welcome message
	fmt.Println()
	fmt.Println("🔑 \033[1;32mNew wallet created!\033[0m")
	fmt.Printf("   Address: \033[1;36m%s\033[0m\n", address)
	fmt.Printf("   Config:  \033[90m%s\033[0m\n", configPath)
	fmt.Println()
	fmt.Println("💡 \033[33mTo use Waffle, get free testnet ETH from:\033[0m")
	fmt.Println("   https://www.alchemy.com/faucets/base-sepolia")
	fmt.Println()

	return cfg, true, nil
}

// loadFromFile reads config from a file
func loadFromFile(path string) (*Config, error) {
	cfg := &Config{
		PricingRate: 1,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// Load reads configuration from ~/.waffle/config.yaml and overrides with environment variables
func Load() (*Config, error) {
	// Ensure config exists (auto-create if needed)
	cfg, _, err := EnsureConfig()
	if err != nil {
		return nil, err
	}

	// Override with Environment Variables
	if val := os.Getenv("RPC_URL"); val != "" {
		cfg.RPCURL = val
	}
	if val := os.Getenv("PRIVATE_KEY"); val != "" {
		cfg.PrivateKey = val
	}
	if val := os.Getenv("SYRUP_TOKEN"); val != "" {
		cfg.SyrupToken = val
	}
	if val := os.Getenv("BAKE_REGISTRY"); val != "" {
		cfg.BakeRegistry = val
	}

	return cfg, nil
}

// GetAddress derives the Ethereum address from the private key
func (c *Config) GetAddress() string {
	if c.PrivateKey == "" {
		return ""
	}

	privateKey, err := crypto.HexToECDSA(c.PrivateKey)
	if err != nil {
		return ""
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return ""
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
}
