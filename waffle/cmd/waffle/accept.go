package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/contracts"
	"github.com/waffle-studio/waffle/internal/wallet"
)

func init() {
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

		acceptSolution(requestID)
	},
}

func acceptSolution(requestID int64) {
	fmt.Println()
	fmt.Printf("%s⏳ Accepting solution for request #%d...%s\n", ColorBlue, requestID, ColorReset)

	// Load config
	privateKey := os.Getenv("PRIVATE_KEY")
	rpcURL := os.Getenv("RPC_URL")
	registryAddr := os.Getenv("BAKE_REGISTRY")

	if privateKey == "" || registryAddr == "" {
		fmt.Printf("%s❌ Missing environment variables%s\n", "\033[31m", ColorReset)
		return
	}

	if rpcURL == "" {
		rpcURL = "http://127.0.0.1:8545"
	}

	// Connect
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		fmt.Printf("%s❌ Failed to connect: %s%s\n", "\033[31m", err, ColorReset)
		return
	}
	defer client.Close()

	registry, err := contracts.NewBakeRegistry(registryAddr, client, privateKey)
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

	// Accept
	receipt, err := registry.AcceptSolution(ctx, big.NewInt(requestID))
	if err != nil {
		fmt.Printf("%s❌ Accept failed: %s%s\n", "\033[31m", err, ColorReset)
		return
	}

	// Format reward
	rewardFloat := wallet.FormatTokenBalance(req.Reward, 18)

	fmt.Println()
	fmt.Printf("%s╔════════════════════════════════════════════════════╗%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s║%s %s✅ SOLUTION ACCEPTED%s                                %s║%s\n", ColorBlue, ColorReset, "\033[32m"+ColorBold, ColorReset, ColorBlue, ColorReset)
	fmt.Printf("%s╠════════════════════════════════════════════════════╣%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s║%s   Request ID: #%-37d%s║%s\n", ColorBlue, ColorReset, requestID, ColorBlue, ColorReset)
	fmt.Printf("%s║%s   🍯 Baker paid: %s%.2f SYRUP%s                        %s║%s\n", ColorBlue, ColorReset, ColorAmber, rewardFloat, ColorReset, ColorBlue, ColorReset)
	fmt.Printf("%s║%s   🔗 Tx: %-42s %s║%s\n", ColorBlue, ColorReset, receipt.TxHash.Hex()[:42], ColorBlue, ColorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════╝%s\n", ColorBlue, ColorReset)
	fmt.Println()
}
