package main

import (
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/cli"
	"github.com/waffle-studio/waffle/internal/contracts"
	"github.com/waffle-studio/waffle/internal/ui"
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
	cli.PrintLoading("Fetching requests...")

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

	// Get next request ID to know how many requests exist
	nextID, err := ctx.Registry.GetNextRequestID(timeoutCtx)
	if err != nil {
		cli.PrintErrorf("Failed to get request count: %s", err)
		return
	}

	const boxWidth = 68

	if nextID.Cmp(big.NewInt(0)) == 0 {
		ui.RequestsHeader("MY BAKE REQUESTS", 54)
		ui.RequestsRow(54, "  No requests yet. Use 'waffle bake' to create one.")
		ui.RequestsFooter(54)
		return
	}

	// Get wallet address for filtering
	myAddress := ctx.Wallet.Address()

	// Display header
	ui.RequestsHeader("MY BAKE REQUESTS", boxWidth)

	// Iterate through all requests
	foundCount := 0
	for i := int64(0); i < nextID.Int64(); i++ {
		req, err := ctx.Registry.GetRequest(timeoutCtx, big.NewInt(i))
		if err != nil {
			continue
		}

		// Only show requests from this wallet
		if req.Requester != myAddress {
			continue
		}

		foundCount++

		// Format status using ui package
		statusIcon, statusText := ui.FormatStatusWithIcon(req.Status)

		// Format reward
		rewardFloat := wallet.FormatTokenBalance(req.Reward, 18)

		// Format code hash (truncated)
		codeHashHex := fmt.Sprintf("0x%x", req.CodeHash[:4])

		content := fmt.Sprintf("%s #%-3d  %s  %s%8.2f SYRUP%s  Hash: %s",
			statusIcon, i, statusText,
			ui.Amber, rewardFloat, ui.Reset,
			codeHashHex)
		ui.RequestsRow(boxWidth, content)
	}

	if foundCount == 0 {
		ui.RequestsRow(boxWidth, "  No requests found for your address.")
	}

	ui.RequestsFooter(boxWidth)
}

// mapStatusToUI converts contract status to ui status
func mapStatusToUI(status uint8) uint8 {
	switch status {
	case contracts.StatusPending:
		return ui.StatusPending
	case contracts.StatusSubmitted:
		return ui.StatusSubmitted
	case contracts.StatusAccepted:
		return ui.StatusAccepted
	case contracts.StatusRejected:
		return ui.StatusRejected
	case contracts.StatusCancelled:
		return ui.StatusCancelled
	default:
		return ui.StatusPending
	}
}
