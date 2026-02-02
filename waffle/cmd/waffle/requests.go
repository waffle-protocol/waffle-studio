package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/config"
	"github.com/waffle-studio/waffle/internal/contracts"
	"github.com/waffle-studio/waffle/internal/wallet"
)

func init() {
	rootCmd.AddCommand(requestsCmd)
}

var requestsCmd = &cobra.Command{
	Use:   "requests",
	Short: "List your bake requests",
	Long:  `Display all bake requests you have created on the BakeRegistry.`,
	Run: func(cmd *cobra.Command, args []string) {
		displayRequests()
	},
}

func displayRequests() {
	fmt.Printf("%s⏳ Fetching requests...%s\n", ColorBlue, ColorReset)

	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("%s❌ Failed to load config: %s%s\n", "\033[31m", err, ColorReset)
		return
	}

	if cfg.PrivateKey == "" || cfg.BakeRegistry == "" {
		fmt.Printf("%s❌ Missing configuration%s\n", "\033[31m", ColorReset)
		fmt.Println("  Please set PRIVATE_KEY and BAKE_REGISTRY in ~/.waffle/config.yaml or env vars.")
		return
	}

	// Connect to RPC
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		fmt.Printf("%s❌ Failed to connect: %s%s\n", "\033[31m", err, ColorReset)
		return
	}
	defer client.Close()

	// Create registry instance
	registry, err := contracts.NewRegistryClient(cfg.BakeRegistry, client, cfg.PrivateKey)
	if err != nil {
		fmt.Printf("%s❌ Failed to setup registry: %s%s\n", "\033[31m", err, ColorReset)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get next request ID to know how many requests exist
	nextID, err := registry.GetNextRequestID(ctx)
	if err != nil {
		fmt.Printf("%s❌ Failed to get request count: %s%s\n", "\033[31m", err, ColorReset)
		return
	}

	if nextID.Cmp(big.NewInt(0)) == 0 {
		fmt.Println()
		fmt.Printf("%s╔════════════════════════════════════════════════════╗%s\n", ColorBlue, ColorReset)
		fmt.Printf("%s║%s %sMY BAKE REQUESTS%s                                    %s║%s\n", ColorBlue, ColorReset, ColorAmber+ColorBold, ColorReset, ColorBlue, ColorReset)
		fmt.Printf("%s╠════════════════════════════════════════════════════╣%s\n", ColorBlue, ColorReset)
		fmt.Printf("%s║%s   No requests yet. Use 'waffle bake' to create one. %s║%s\n", ColorBlue, ColorReset, ColorBlue, ColorReset)
		fmt.Printf("%s╚════════════════════════════════════════════════════╝%s\n", ColorBlue, ColorReset)
		fmt.Println()
		return
	}

	// Get wallet address for filtering
	w, err := wallet.New(cfg)
	if err != nil {
		fmt.Printf("%s❌ Failed to create wallet: %s%s\n", "\033[31m", err, ColorReset)
		return
	}
	myAddress := w.Address()

	// Display header
	fmt.Println()
	fmt.Printf("%s╔════════════════════════════════════════════════════════════════╗%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s║%s %sMY BAKE REQUESTS%s                                                %s║%s\n", ColorBlue, ColorReset, ColorAmber+ColorBold, ColorReset, ColorBlue, ColorReset)
	fmt.Printf("%s╠════════════════════════════════════════════════════════════════╣%s\n", ColorBlue, ColorReset)

	// Iterate through all requests
	foundCount := 0
	for i := int64(0); i < nextID.Int64(); i++ {
		req, err := registry.GetRequest(ctx, big.NewInt(i))
		if err != nil {
			continue
		}

		// Only show requests from this wallet
		if req.Requester != myAddress {
			continue
		}

		foundCount++

		// Format status
		var statusIcon, statusText string
		switch req.Status {
		case contracts.StatusPending:
			statusIcon = "⏳"
			statusText = "\033[33mPENDING\033[0m  "
		case contracts.StatusSubmitted:
			statusIcon = "📝"
			statusText = "\033[36mSUBMITTED\033[0m"
		case contracts.StatusAccepted:
			statusIcon = "✅"
			statusText = "\033[32mACCEPTED\033[0m "
		case contracts.StatusRejected:
			statusIcon = "❌"
			statusText = "\033[31mREJECTED\033[0m "
		case contracts.StatusCancelled:
			statusIcon = "🚫"
			statusText = "\033[90mCANCELLED\033[0m"
		}

		// Format reward
		rewardFloat := wallet.FormatTokenBalance(req.Reward, 18)

		// Format code hash (truncated)
		codeHashHex := fmt.Sprintf("0x%x", req.CodeHash[:4])

		fmt.Printf("%s║%s %s #%-3d  %s  %s%8.2f SYRUP%s  Hash: %s   %s║%s\n",
			ColorBlue, ColorReset,
			statusIcon, i, statusText,
			ColorAmber, rewardFloat, ColorReset,
			codeHashHex,
			ColorBlue, ColorReset)
	}

	if foundCount == 0 {
		fmt.Printf("%s║%s   No requests found for your address.                          %s║%s\n", ColorBlue, ColorReset, ColorBlue, ColorReset)
	}

	fmt.Printf("%s╚════════════════════════════════════════════════════════════════╝%s\n", ColorBlue, ColorReset)
	fmt.Println()
}
