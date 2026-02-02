package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// ANSI color codes
const (
	ColorAmber = "\033[38;5;214m"
	ColorBlue  = "\033[34m"
	ColorBold  = "\033[1m"
	ColorReset = "\033[0m"
)

var rootCmd = &cobra.Command{
	Use:   "waffle",
	Short: "Waffle - P2P LLM Token Marketplace CLI",
	Long:  `Waffle is a decentralized marketplace for LLM tokens, built on blockchain technology.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Server Rack style ASCII art
		fmt.Println()
		fmt.Printf("%s[ 🧇🧇🧇 ]%s  %sWAFFLE STACK%s online.\n", ColorAmber, ColorReset, ColorBold, ColorReset)
		fmt.Printf("%s[ 🧇🍯🧇 ]%s  %sConnecting to P2P Grid...%s\n", ColorAmber, ColorReset, ColorBlue, ColorReset)
		fmt.Printf("%s[ 🧇🧇🧇 ]%s  v0.1.0 ready.\n", ColorAmber, ColorReset)
		fmt.Println()
		fmt.Printf("%s🥞 Waffle 스택이 활성화되었습니다.%s\n", ColorBold, ColorReset)
		fmt.Println()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
