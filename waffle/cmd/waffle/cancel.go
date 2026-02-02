package main

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/cli"
	"github.com/waffle-studio/waffle/internal/contracts"
	"github.com/waffle-studio/waffle/internal/ui"
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
			cli.PrintErrorf("Invalid request ID: %s", args[0])
			return
		}

		cancelRequest(requestID)
	},
}

func cancelRequest(requestID int64) {
	fmt.Println()
	cli.PrintLoading(fmt.Sprintf("Cancelling request #%d...", requestID))

	// Bootstrap with registry
	ctx, err := cli.BootstrapWithRegistry()
	if err != nil {
		cli.PrintError(err.Error())
		cli.PrintConfigHint()
		return
	}
	defer ctx.Close()

	timeoutCtx, cancel := cli.WithTimeout()
	defer cancel()

	// Check request status first
	req, err := ctx.Registry.GetRequest(timeoutCtx, big.NewInt(requestID))
	if err != nil {
		cli.PrintErrorf("Failed to get request: %s", err)
		return
	}

	if req.Status != contracts.StatusPending {
		cli.PrintErrorf("Can only cancel PENDING requests (current status: %d)", req.Status)
		return
	}

	// Cancel
	receipt, err := ctx.Registry.CancelRequest(timeoutCtx, big.NewInt(requestID))
	if err != nil {
		cli.PrintErrorf("Cancel failed: %s", err)
		return
	}

	// Display result box
	box := ui.NewBox()
	fmt.Println()
	box.Header("🚫 REQUEST CANCELLED")
	box.Row(fmt.Sprintf("   Request ID: #%-37d", requestID))
	box.Row("   💰 SYRUP refunded to your wallet")
	box.Row(fmt.Sprintf("   🔗 Tx: %-42s", receipt.TxHash.Hex()[:42]))
	box.Footer()
	fmt.Println()
}
