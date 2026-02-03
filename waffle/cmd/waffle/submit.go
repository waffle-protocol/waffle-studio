package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/waffle-studio/waffle/internal/cli"
	"github.com/waffle-studio/waffle/internal/config"
	"github.com/waffle-studio/waffle/internal/contracts"
	"github.com/waffle-studio/waffle/internal/ui"
	"github.com/waffle-studio/waffle/internal/wallet"
)

var (
	tokenUsage   uint64
	solutionText string
	useRelay     bool
)

func init() {
	submitCmd.Flags().Uint64VarP(&tokenUsage, "usage", "u", 0, "Token usage for the solution")
	submitCmd.Flags().StringVarP(&solutionText, "solution", "s", "", "Solution text (will be auto-hashed)")
	submitCmd.Flags().BoolVar(&useRelay, "relay", true, "Use gasless relay (default: true)")
	rootCmd.AddCommand(submitCmd)
}

var submitCmd = &cobra.Command{
	Use:   "submit <request_id>",
	Short: "Submit a solution for a pending request (Baker)",
	Long: `As a Baker, submit a solution for a PENDING bake request to earn SYRUP reward.

Examples:
  waffle submit 0 --solution "Here is my solution text"
  waffle submit 0 -s "Solution content" --usage 100
  waffle submit 0 -s "Solution" --relay=false  # Direct TX (requires ETH)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		requestID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			cli.PrintErrorf("Invalid request ID: %s", args[0])
			return
		}

		if solutionText == "" {
			cli.PrintError("Please provide a solution with --solution or -s flag")
			fmt.Println()
			fmt.Println("Example: waffle submit 0 --solution 'Your solution here'")
			return
		}

		if tokenUsage == 0 {
			cli.PrintWarning("Token usage not specified (using 0). Use --usage flag.")
		}

		// Auto-hash the solution text
		hash := sha256.Sum256([]byte(solutionText))
		solutionHashStr := hex.EncodeToString(hash[:])

		submitSolution(requestID, solutionHashStr, tokenUsage, solutionText, useRelay)
	},
}

// RelayRequest for gasless submission
type RelayRequest struct {
	RequestID    uint64 `json:"request_id"`
	SolutionHash string `json:"solution_hash"`
	TokenUsage   uint64 `json:"token_usage"`
	Signature    string `json:"signature"`
	Signer       string `json:"signer"`
}

// RelayResponse from relay server
type RelayResponse struct {
	Success bool   `json:"success"`
	TxHash  string `json:"tx_hash"`
	Error   string `json:"error"`
}

func submitViaRelay(cfg *config.Config, requestID int64, solutionHash string, usage uint64, signerAddress string) (string, error) {
	relayURL := cfg.RelayURL
	if relayURL == "" {
		relayURL = config.DefaultRelayURL
	}

	// 1. Create message to sign
	// Format matches Solidity: "SubmitSolution(requestId,solutionHash,tokenUsage)"
	// Note: Solidity's toString(bytes32) returns raw bytes, not hex string
	solutionHashBytes, err := hex.DecodeString(solutionHash)
	if err != nil {
		return "", fmt.Errorf("invalid solution hash: %w", err)
	}
	msg := fmt.Sprintf("SubmitSolution(%d,%s,%d)", requestID, string(solutionHashBytes), usage)

	// This hash matches 'messageHash' in Solidity
	msgHash := crypto.Keccak256Hash([]byte(msg))

	// Create Ethereum Signed Message Hash (EIP-191)
	// Matches: keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", messageHash))
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msgHash))
	prefixedHash := crypto.Keccak256Hash(append([]byte(prefix), msgHash.Bytes()...))

	// 2. Sign the prefixed hash
	keyHex := strings.TrimPrefix(cfg.PrivateKey, "0x")
	privateKey, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}

	signatureBytes, err := crypto.Sign(prefixedHash.Bytes(), privateKey)
	if err != nil {
		return "", fmt.Errorf("sign failed: %w", err)
	}
	// Adjust V value for Ethereum (27 or 28)
	if len(signatureBytes) == 65 {
		signatureBytes[64] += 27
	}
	signature := hex.EncodeToString(signatureBytes)

	// 3. Send request to relay
	reqBody := RelayRequest{
		RequestID:    uint64(requestID),
		SolutionHash: solutionHash,
		TokenUsage:   usage,
		Signature:    signature,
		Signer:       signerAddress,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Post(relayURL+"/submit-solution", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("relay request failed: %w", err)
	}
	defer resp.Body.Close()

	var relayResp RelayResponse
	if err := json.NewDecoder(resp.Body).Decode(&relayResp); err != nil {
		return "", fmt.Errorf("failed to parse relay response: %w", err)
	}

	if !relayResp.Success {
		return "", fmt.Errorf("relay reported error: %s", relayResp.Error)
	}

	return relayResp.TxHash, nil
}

func submitSolution(requestID int64, solutionHashStr string, usage uint64, originalText string, tryRelay bool) {
	fmt.Println()
	cli.PrintLoading(fmt.Sprintf("Submitting solution for request #%d...", requestID))

	// Load config
	cfg, err := config.Load()
	if err != nil {
		cli.PrintError(err.Error())
		return
	}

	// Bootstrap with registry for status check
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

	var txHash string

	// Try relay first (gasless)
	if tryRelay && cfg.RelayURL != "" {
		fmt.Printf("  %s Using gasless relay...%s\n", ui.TextAmber.Render("⏳"), ui.Reset)

		// Get signer address
		w, err := wallet.New(cfg)
		if err == nil {
			txHash, err = submitViaRelay(cfg, requestID, solutionHashStr, usage, w.AddressHex())
			w.Close()
		}

		if err != nil {
			fmt.Printf("  ⚠️  Relay failed, trying direct TX...\n")
			tryRelay = false
		}
	}

	// Fallback to direct TX
	if !tryRelay || txHash == "" {
		// Parse solution hash
		var solutionHash [32]byte
		hashBytes, err := hex.DecodeString(solutionHashStr)
		if err != nil || len(hashBytes) != 32 {
			cli.PrintError("Invalid solution hash")
			return
		}
		copy(solutionHash[:], hashBytes)

		// Submit solution directly
		receipt, err := ctx.Registry.SubmitSolution(timeoutCtx, big.NewInt(requestID), solutionHash, big.NewInt(int64(usage)))
		if err != nil {
			cli.PrintErrorf("Submit failed: %s", err)
			return
		}
		txHash = receipt.TxHash.Hex()
	}

	// Format reward
	rewardFloat := wallet.FormatTokenBalance(req.Reward, 18)

	// Display result box
	box := ui.NewBox()
	fmt.Println()
	box.Header("🍳 SOLUTION SUBMITTED")
	box.Row(fmt.Sprintf("   Request ID: #%-37d", requestID))
	if len(originalText) > 40 {
		box.Row(fmt.Sprintf("   📝 Solution: %.40s...", originalText))
	} else {
		box.Row(fmt.Sprintf("   📝 Solution: %-40s", originalText))
	}
	box.Row(fmt.Sprintf("   🍯 Potential reward: %s%.2f SYRUP%s", ui.Amber, rewardFloat, ui.Reset))
	if len(txHash) > 42 {
		box.Row(fmt.Sprintf("   🔗 Tx: %-42s", txHash[:42]))
	} else {
		box.Row(fmt.Sprintf("   🔗 Tx: %-42s", txHash))
	}
	box.Footer()
	fmt.Println()
	fmt.Printf("%s  ⏳ Waiting for requester to accept your solution...%s\n", ui.Blue, ui.Reset)
	fmt.Println()
}
