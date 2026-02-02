package main

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/cli"
	"github.com/waffle-studio/waffle/internal/contracts"
	"github.com/waffle-studio/waffle/internal/ui"
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
			cli.PrintErrorf("Invalid request ID: %s", args[0])
			return
		}

		acceptSolution(requestID, paymentAmount)
	},
}

func acceptSolution(requestID int64, amount float64) {
	fmt.Println()
	cli.PrintLoading(fmt.Sprintf("Accepting solution for request #%d...", requestID))

	// Bootstrap with full context
	ctx, err := cli.BootstrapFull()
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

	if req.Status != contracts.StatusSubmitted {
		cli.PrintErrorf("Can only accept SUBMITTED solutions (current status: %d)", req.Status)
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
		rate := ctx.Config.PricingRate
		if rate == 0 {
			rate = 1 // Default pricing rate
		}
		calculatedAmount := float64(usage) * float64(rate)

		fmt.Printf("  %s Baker used %d tokens (Rate: %.2f)\n", ui.TextAmber.Render("📊"), usage, float64(rate))
		fmt.Printf("  %s Suggested payment: %.2f SYRUP\n", ui.TextAmber.Render("💰"), calculatedAmount)

		paymentWei = wallet.ToWei(calculatedAmount, 18)
	}

	// Ensure payment doesn't exceed reward cap
	if paymentWei.Cmp(req.Reward) > 0 {
		cli.PrintWarningf("Payment exceeds escrowed reward (%.2f). Capping at reward.", wallet.FormatTokenBalance(req.Reward, 18))
		paymentWei = req.Reward
	}

	// Accept
	receipt, err := ctx.Registry.AcceptSolution(timeoutCtx, big.NewInt(requestID), paymentWei)
	if err != nil {
		cli.PrintErrorf("Accept failed: %s", err)
		return
	}

	// Format amounts
	rewardFloat := wallet.FormatTokenBalance(paymentWei, 18)
	refundWei := new(big.Int).Sub(req.Reward, paymentWei)
	refundFloat := wallet.FormatTokenBalance(refundWei, 18)

	// Display result box
	box := ui.NewBox()
	fmt.Println()
	box.Header("✅ SOLUTION ACCEPTED")
	box.Row(fmt.Sprintf("   Request ID: #%-37d", requestID))
	box.Row(fmt.Sprintf("   🍯 Baker paid: %s%.2f SYRUP%s", ui.Amber, rewardFloat, ui.Reset))
	if refundWei.Cmp(big.NewInt(0)) > 0 {
		box.Row(fmt.Sprintf("   ↩️  Refunded:   %s%.2f SYRUP%s", ui.Amber, refundFloat, ui.Reset))
	}
	box.Row(fmt.Sprintf("   🔗 Tx: %-42s", receipt.TxHash.Hex()[:42]))
	box.Footer()
	fmt.Println()
}
