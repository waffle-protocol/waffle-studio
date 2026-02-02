package main

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/config"
	"github.com/waffle-studio/waffle/internal/contracts"
	"github.com/waffle-studio/waffle/internal/wallet"
)

var paymentAmount float64

func init() {
	acceptCmd.Flags().Float64VarP(&paymentAmount, "amount", "a", 0, "Amount of SYRUP to pay (overrides calculated price)")
	rootCmd.AddCommand(acceptCmd)
}

var acceptCmd = &cobra.Command{
	Use:   "accept <request_id>",
	Short: "Accept a submitted solution",
	Long:  `Accept a SUBMITTED solution and pay the baker their SYRUP reward.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		requestID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Printf("%s❌ Invalid request ID: %s%s\n", "\033[31m", args[0], ColorReset)
			return
		}

		acceptSolution(requestID, paymentAmount)
	},
}

func acceptSolution(requestID int64, amount float64) {
	fmt.Println()
	fmt.Printf("%s⏳ Accepting solution for request #%d...%s\n", ColorBlue, requestID, ColorReset)

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

	// Connect
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		fmt.Printf("%s❌ Failed to connect: %s%s\n", "\033[31m", err, ColorReset)
		return
	}
	defer client.Close()

	registry, err := contracts.NewRegistryClient(cfg.BakeRegistry, client, cfg.PrivateKey)
	if err != nil {
		fmt.Printf("%s❌ Failed to setup registry: %s%s\n", "\033[31m", err, ColorReset)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Check request status first
	req, err := registry.GetRequest(ctx, big.NewInt(requestID))
	if err != nil {
		fmt.Printf("%s❌ Failed to get request: %s%s\n", "\033[31m", err, ColorReset)
		return
	}

	if req.Status != contracts.StatusSubmitted {
		fmt.Printf("%s❌ Can only accept SUBMITTED solutions (current status: %d)%s\n", "\033[31m", req.Status, ColorReset)
		return
	}

	// Calculate payment amount
	var paymentWei *big.Int

	if amount > 0 {
		// Use manually specified amount
		paymentWei = wallet.ToWei(amount, 18)
	} else {
		// Calculate based on usage
		usage := req.TokenUsage.Uint64()
		rate := cfg.PricingRate
		if rate == 0 {
			rate = 1 // Default pricing rate
		}
		calculatedAmount := float64(usage) * float64(rate)

		fmt.Printf("  %s Baker used %d tokens (Rate: %.2f)\n", textAmber.Render("📊"), usage, float64(rate))
		fmt.Printf("  %s Suggested payment: %.2f SYRUP\n", textAmber.Render("💰"), calculatedAmount)

		// Auto-accept calculated amount for now (prompting needs TUI or bufio)
		// For CLI UX improvement, we just use calculated amount if not specified
		paymentWei = wallet.ToWei(calculatedAmount, 18)
	}

	// Ensure payment doesn't exceed reward cap
	if paymentWei.Cmp(req.Reward) > 0 {
		fmt.Printf("%s⚠️  Payment exceeds escrowed reward (%.2f). Capping at reward.%s\n", "\033[33m", wallet.FormatTokenBalance(req.Reward, 18), ColorReset)
		paymentWei = req.Reward
	}

	// Accept
	receipt, err := registry.AcceptSolution(ctx, big.NewInt(requestID), paymentWei)
	if err != nil {
		fmt.Printf("%s❌ Accept failed: %s\n", "\033[31m", err)
		return
	}

	// Format reward
	rewardFloat := wallet.FormatTokenBalance(paymentWei, 18)
	refundWei := new(big.Int).Sub(req.Reward, paymentWei)
	refundFloat := wallet.FormatTokenBalance(refundWei, 18)

	fmt.Println()
	fmt.Printf("%s╔════════════════════════════════════════════════════╗%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s║%s %s✅ SOLUTION ACCEPTED%s                                %s║%s\n", ColorBlue, ColorReset, "\033[32m"+ColorBold, ColorReset, ColorBlue, ColorReset)
	fmt.Printf("%s╠════════════════════════════════════════════════════╣%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s║%s   Request ID: #%-37d%s║%s\n", ColorBlue, ColorReset, requestID, ColorBlue, ColorReset)
	fmt.Printf("%s║%s   🍯 Baker paid: %s%.2f SYRUP%s                        %s║%s\n", ColorBlue, ColorReset, ColorAmber, rewardFloat, ColorReset, ColorBlue, ColorReset)
	if refundWei.Cmp(big.NewInt(0)) > 0 {
		fmt.Printf("%s║%s   ↩️  Refunded:   %s%.2f SYRUP%s                        %s║%s\n", ColorBlue, ColorReset, ColorAmber, refundFloat, ColorReset, ColorBlue, ColorReset)
	}
	fmt.Printf("%s║%s   🔗 Tx: %-42s %s║%s\n", ColorBlue, ColorReset, receipt.TxHash.Hex()[:42], ColorBlue, ColorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════╝%s\n", ColorBlue, ColorReset)
	fmt.Println()
}
