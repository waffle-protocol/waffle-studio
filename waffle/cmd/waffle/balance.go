package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/cli"
	"github.com/waffle-studio/waffle/internal/ui"
	"github.com/waffle-studio/waffle/internal/wallet"
)

func init() {
	rootCmd.AddCommand(balanceCmd)
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Check wallet balance",
	Long:  `Display the current SYRUP token balance for the connected wallet.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.PrintLoading("Connecting to blockchain...")

		// Bootstrap with wallet
		ctx, err := cli.BootstrapWithWallet()
		if err != nil {
			cli.PrintError(err.Error())
			cli.PrintConfigHint()
			return
		}
		defer ctx.Close()

		if ctx.Config.SyrupToken == "" {
			cli.PrintError("SYRUP_TOKEN not configured")
			cli.PrintConfigHint()
			return
		}

		// Query token balance
		timeoutCtx, cancel := cli.WithQuickTimeout()
		defer cancel()

		tokenAddr := common.HexToAddress(ctx.Config.SyrupToken)
		balance, err := ctx.Wallet.GetTokenBalance(timeoutCtx, tokenAddr)
		if err != nil {
			cli.PrintErrorf("Failed to get balance: %s", err.Error())
			return
		}

		// Format and display
		balanceFormatted := wallet.FormatTokenBalance(balance, 18)
		address := ctx.Wallet.AddressHex()

		ui.WalletBox(address, balanceFormatted, "SYRUP")
	},
}
