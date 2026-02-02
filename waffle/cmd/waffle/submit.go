package main

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/cli"
	"github.com/waffle-studio/waffle/internal/contracts"
	"github.com/waffle-studio/waffle/internal/ui"
	"github.com/waffle-studio/waffle/internal/wallet"
)

var tokenUsage uint64

func init() {
	submitCmd.Flags().Uint64VarP(&tokenUsage, "usage", "u", 0, "Token usage for the solution")
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
			cli.PrintErrorf("Invalid request ID: %s", args[0])
			return
		}

		solutionHashStr := args[1]
		if !strings.HasPrefix(solutionHashStr, "0x") {
			solutionHashStr = "0x" + solutionHashStr
		}

		if tokenUsage == 0 {
			cli.PrintWarning("Token usage not specified (using 0). Use --usage flag.")
		}

		submitSolution(requestID, solutionHashStr, tokenUsage)
	},
}

func submitSolution(requestID int64, solutionHashStr string, usage uint64) {
	fmt.Println()
	cli.PrintLoading(fmt.Sprintf("Submitting solution for request #%d...", requestID))

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
		cli.PrintErrorf("Can only submit solutions for PENDING requests (current status: %d)", req.Status)
		return
	}

	// Parse solution hash
	solutionHashStr = strings.TrimPrefix(solutionHashStr, "0x")
	if len(solutionHashStr) != 64 {
		cli.PrintError("Invalid solution hash length (expected 64 hex chars)")
		return
	}

	var solutionHash [32]byte
	for i := 0; i < 32; i++ {
		b, err := strconv.ParseUint(solutionHashStr[i*2:i*2+2], 16, 8)
		if err != nil {
			cli.PrintErrorf("Invalid solution hash format: %s", err)
			return
		}
		solutionHash[i] = byte(b)
	}

	// Submit solution
	receipt, err := ctx.Registry.SubmitSolution(timeoutCtx, big.NewInt(requestID), solutionHash, big.NewInt(int64(usage)))
	if err != nil {
		cli.PrintErrorf("Submit failed: %s", err)
		return
	}

	// Format reward
	rewardFloat := wallet.FormatTokenBalance(req.Reward, 18)

	// Display result box
	box := ui.NewBox()
	fmt.Println()
	box.Header("🍳 SOLUTION SUBMITTED")
	box.Row(fmt.Sprintf("   Request ID: #%-37d", requestID))
	box.Row(fmt.Sprintf("   🍯 Potential reward: %s%.2f SYRUP%s", ui.Amber, rewardFloat, ui.Reset))
	box.Row(fmt.Sprintf("   🔗 Tx: %-42s", receipt.TxHash.Hex()[:42]))
	box.Footer()
	fmt.Println()
	fmt.Printf("%s  ⏳ Waiting for requester to accept your solution...%s\n", ui.Blue, ui.Reset)
	fmt.Println()
}
