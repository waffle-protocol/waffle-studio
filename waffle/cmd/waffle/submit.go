package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/contracts"
	"github.com/waffle-studio/waffle/internal/wallet"
)

func init() {
	rootCmd.AddCommand(submitCmd)
}

var submitCmd = &cobra.Command{
	Use:   "submit <request_id> <solution_hash>",
	Short: "Submit a solution for a pending request (Baker)",
	Long:  `As a Baker, submit a solution hash for a PENDING bake request to earn SYRUP reward.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		requestID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Printf("%s❌ Invalid request ID: %s%s\n", "\033[31m", args[0], ColorReset)
			return
		}

		solutionHashStr := args[1]
		if !strings.HasPrefix(solutionHashStr, "0x") {
			solutionHashStr = "0x" + solutionHashStr
		}

		submitSolution(requestID, solutionHashStr)
	},
}

func submitSolution(requestID int64, solutionHashStr string) {
	fmt.Println()
	fmt.Printf("%s⏳ Submitting solution for request #%d...%s\n", ColorBlue, requestID, ColorReset)

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

	if req.Status != contracts.StatusPending {
		fmt.Printf("%s❌ Can only submit solutions for PENDING requests (current status: %d)%s\n", "\033[31m", req.Status, ColorReset)
		return
	}

	// Parse solution hash
	solutionHashStr = strings.TrimPrefix(solutionHashStr, "0x")
	if len(solutionHashStr) != 64 {
		fmt.Printf("%s❌ Invalid solution hash length (expected 64 hex chars)%s\n", "\033[31m", ColorReset)
		return
	}

	var solutionHash [32]byte
	for i := 0; i < 32; i++ {
		b, err := strconv.ParseUint(solutionHashStr[i*2:i*2+2], 16, 8)
		if err != nil {
			fmt.Printf("%s❌ Invalid solution hash format: %s%s\n", "\033[31m", err, ColorReset)
			return
		}
		solutionHash[i] = byte(b)
	}

	// Submit solution
	receipt, err := registry.SubmitSolution(ctx, big.NewInt(requestID), solutionHash)
	if err != nil {
		fmt.Printf("%s❌ Submit failed: %s%s\n", "\033[31m", err, ColorReset)
		return
	}

	// Format reward
	rewardFloat := wallet.FormatTokenBalance(req.Reward, 18)

	fmt.Println()
	fmt.Printf("%s╔════════════════════════════════════════════════════╗%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s║%s %s🍳 SOLUTION SUBMITTED%s                               %s║%s\n", ColorBlue, ColorReset, ColorAmber+ColorBold, ColorReset, ColorBlue, ColorReset)
	fmt.Printf("%s╠════════════════════════════════════════════════════╣%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s║%s   Request ID: #%-37d%s║%s\n", ColorBlue, ColorReset, requestID, ColorBlue, ColorReset)
	fmt.Printf("%s║%s   🍯 Potential reward: %s%.2f SYRUP%s                  %s║%s\n", ColorBlue, ColorReset, ColorAmber, rewardFloat, ColorReset, ColorBlue, ColorReset)
	fmt.Printf("%s║%s   🔗 Tx: %-42s %s║%s\n", ColorBlue, ColorReset, receipt.TxHash.Hex()[:42], ColorBlue, ColorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════╝%s\n", ColorBlue, ColorReset)
	fmt.Println()
	fmt.Printf("%s  ⏳ Waiting for requester to accept your solution...%s\n", ColorBlue, ColorReset)
	fmt.Println()
}
