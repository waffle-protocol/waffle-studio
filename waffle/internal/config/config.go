package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RPCURL       string `yaml:"rpc_url"`
	PrivateKey   string `yaml:"private_key"`
	SyrupToken   string `yaml:"syrup_token"`
	BakeRegistry string `yaml:"bake_registry"`
	PricingRate  uint64 `yaml:"pricing_rate"` // Rate per character+usage
}

// Load reads configuration from ~/.waffle/config.yaml and overrides with environment variables
func Load() (*Config, error) {
	cfg := &Config{
		PricingRate: 1, // Default pricing rate
	}

	// 1. Load from file
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(homeDir, ".waffle", "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			data, err := os.ReadFile(configPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("failed to parse config file: %w", err)
			}
		}
	}

	// 2. Override with Environment Variables
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
	// PricingRate likely won't be in env var often, but could add if needed

	// 3. Defaults for testnet (if not set)
	// TODO: Uncomment this after testnet deployment is stable
	// if cfg.RPCURL == "" {
	// 	cfg.RPCURL = "https://sepolia.base.org"
	// }

	return cfg, nil
}
