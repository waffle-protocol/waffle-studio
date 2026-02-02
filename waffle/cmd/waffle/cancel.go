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
)

func init() {
	rootCmd.AddCommand(cancelCmd)
}

var cancelCmd = &cobra.Command{
	Use:   "cancel <request_id>",
	Short: "Cancel a pending bake request",
	Long:  `Cancel a PENDING bake request and receive your SYRUP refund.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		requestID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Printf("%s❌ Invalid request ID: %s%s\n", "\033[31m", args[0], ColorReset)
			return
		}

		cancelRequest(requestID)
	},
}

func cancelRequest(requestID int64) {
	fmt.Println()
	fmt.Printf("%s⏳ Cancelling request #%d...%s\n", ColorBlue, requestID, ColorReset)

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
		fmt.Printf("%s❌ Can only cancel PENDING requests (current status: %d)%s\n", "\033[31m", req.Status, ColorReset)
		return
	}

	// Cancel
	receipt, err := registry.CancelRequest(ctx, big.NewInt(requestID))
	if err != nil {
		fmt.Printf("%s❌ Cancel failed: %s%s\n", "\033[31m", err, ColorReset)
		return
	}

	fmt.Println()
	fmt.Printf("%s╔════════════════════════════════════════════════════╗%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s║%s %s🚫 REQUEST CANCELLED%s                                %s║%s\n", ColorBlue, ColorReset, ColorAmber+ColorBold, ColorReset, ColorBlue, ColorReset)
	fmt.Printf("%s╠════════════════════════════════════════════════════╣%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s║%s   Request ID: #%-37d%s║%s\n", ColorBlue, ColorReset, requestID, ColorBlue, ColorReset)
	fmt.Printf("%s║%s   💰 SYRUP refunded to your wallet                 %s║%s\n", ColorBlue, ColorReset, ColorBlue, ColorReset)
	fmt.Printf("%s║%s   🔗 Tx: %-42s %s║%s\n", ColorBlue, ColorReset, receipt.TxHash.Hex()[:42], ColorBlue, ColorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════╝%s\n", ColorBlue, ColorReset)
	fmt.Println()
}
