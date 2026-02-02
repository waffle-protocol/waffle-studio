package cli

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/waffle-studio/waffle/internal/config"
	"github.com/waffle-studio/waffle/internal/contracts"
	"github.com/waffle-studio/waffle/internal/ui"
	"github.com/waffle-studio/waffle/internal/wallet"
)

// ClientContext holds all initialized clients for CLI commands
type ClientContext struct {
	Config   *config.Config
	Client   *ethclient.Client
	Registry *contracts.RegistryClient
	Wallet   *wallet.Wallet
}

// BootstrapOptions specifies which components to initialize
type BootstrapOptions struct {
	NeedWallet   bool
	NeedRegistry bool
	NeedClient   bool
}

// Bootstrap initializes CLI context with requested components
func Bootstrap(opts BootstrapOptions) (*ClientContext, error) {
	ctx := &ClientContext{}

	// Load configuration (auto-creates if needed)
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	ctx.Config = cfg

	// Initialize wallet if needed
	if opts.NeedWallet {
		w, err := wallet.New(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create wallet: %w", err)
		}
		ctx.Wallet = w
		ctx.Client = w.Client()
	}

	// Initialize RPC client if needed (and not already from wallet)
	if opts.NeedClient && ctx.Client == nil {
		client, err := ethclient.Dial(cfg.RPCURL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to RPC: %w", err)
		}
		ctx.Client = client
	}

	// Initialize registry if needed
	if opts.NeedRegistry {
		if cfg.BakeRegistry == "" {
			return nil, fmt.Errorf("BAKE_REGISTRY not configured")
		}
		if ctx.Client == nil {
			client, err := ethclient.Dial(cfg.RPCURL)
			if err != nil {
				return nil, fmt.Errorf("failed to connect to RPC: %w", err)
			}
			ctx.Client = client
		}
		registry, err := contracts.NewRegistryClient(cfg.BakeRegistry, ctx.Client, cfg.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to setup registry: %w", err)
		}
		ctx.Registry = registry
	}

	return ctx, nil
}

// Close releases all resources
func (ctx *ClientContext) Close() {
	// Only close client if we don't have a wallet (wallet manages its own client)
	if ctx.Wallet != nil {
		ctx.Wallet.Close()
	} else if ctx.Client != nil {
		ctx.Client.Close()
	}
}

// MustBootstrap is like Bootstrap but prints error and exits on failure
func MustBootstrap(opts BootstrapOptions) *ClientContext {
	ctx, err := Bootstrap(opts)
	if err != nil {
		PrintError(err.Error())
		PrintConfigHint()
		return nil
	}
	return ctx
}

// BootstrapWithWallet initializes context with wallet
func BootstrapWithWallet() (*ClientContext, error) {
	return Bootstrap(BootstrapOptions{
		NeedWallet: true,
	})
}

// BootstrapWithRegistry initializes context with registry client
func BootstrapWithRegistry() (*ClientContext, error) {
	return Bootstrap(BootstrapOptions{
		NeedRegistry: true,
		NeedClient:   true,
	})
}

// BootstrapFull initializes context with all components
func BootstrapFull() (*ClientContext, error) {
	return Bootstrap(BootstrapOptions{
		NeedWallet:   true,
		NeedRegistry: true,
	})
}

// PrintLoading prints a loading message
func PrintLoading(msg string) {
	fmt.Println(ui.Info(msg))
}
