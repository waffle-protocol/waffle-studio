package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/waffle-studio/waffle/internal/contracts"
)

const (
	testRPC          = "http://127.0.0.1:8545"
	// Anvil Account #0 (Requester)
	testPrivateKey   = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	// Anvil Account #1 (Baker) - different account for submitting solutions
	bakerPrivateKey  = "0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
	testSyrupToken   = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
	testBakeRegistry = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
)

func TestE2EFullFlow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Connect to blockchain
	fmt.Println("=== E2E Test: Full Bake Flow ===")
	fmt.Println()
	fmt.Println("[1/6] Connecting to blockchain...")

	client, err := ethclient.Dial(testRPC)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()
	fmt.Println("      Connected to Anvil!")

	// 2. Initialize registry client
	fmt.Println("[2/6] Initializing registry client...")
	registry, err := contracts.NewRegistryClient(testBakeRegistry, client, testPrivateKey)
	if err != nil {
		t.Fatalf("Failed to create registry client: %v", err)
	}
	fmt.Println("      Registry client ready!")

	// 3. Create a bake request
	fmt.Println("[3/6] Creating bake request...")
	testCode := "function hello() { return 'world'; }"
	codeHash := contracts.HashCode(testCode)
	reward := new(big.Int)
	reward.SetString("10000000000000000000", 10) // 10 SYRUP (10 * 10^18)

	// Approve token first
	fmt.Println("      Approving SYRUP tokens...")
	_, err = registry.ApproveToken(ctx, testSyrupToken, reward)
	if err != nil {
		t.Fatalf("Failed to approve: %v", err)
	}

	// Create request
	fmt.Println("      Sending createRequest transaction...")
	receipt, requestID, err := registry.CreateRequest(ctx, codeHash, reward)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	fmt.Printf("      Request created! ID: %d, Tx: %s\n", requestID, receipt.TxHash.Hex()[:18]+"...")

	// 4. Verify request was created
	fmt.Println("[4/6] Verifying request on-chain...")
	req, err := registry.GetRequest(ctx, requestID)
	if err != nil {
		t.Fatalf("Failed to get request: %v", err)
	}

	if req.Status != contracts.StatusPending {
		t.Fatalf("Expected status Pending(0), got %d", req.Status)
	}
	if req.Reward.Cmp(reward) != 0 {
		t.Fatalf("Reward mismatch: expected %s, got %s", reward.String(), req.Reward.String())
	}
	fmt.Printf("      Status: Pending, Reward: %s wei\n", req.Reward.String())

	// 5. Submit solution (simulating Baker with different account)
	fmt.Println("[5/6] Submitting solution (as Baker - Account #1)...")

	// Create a separate registry client for the Baker
	bakerRegistry, err := contracts.NewRegistryClient(testBakeRegistry, client, bakerPrivateKey)
	if err != nil {
		t.Fatalf("Failed to create baker registry client: %v", err)
	}

	solutionCode := "function hello() { return 'Hello, World!'; }"
	solutionHash := contracts.HashCode(solutionCode)
	tokenUsage := big.NewInt(100) // 100 tokens used

	_, err = bakerRegistry.SubmitSolution(ctx, requestID, solutionHash, tokenUsage)
	if err != nil {
		t.Fatalf("Failed to submit solution: %v", err)
	}

	// Verify status changed
	req, _ = registry.GetRequest(ctx, requestID)
	if req.Status != contracts.StatusSubmitted {
		t.Fatalf("Expected status Submitted(1), got %d", req.Status)
	}
	fmt.Printf("      Solution submitted! Status: Submitted, TokenUsage: %d\n", tokenUsage)

	// 6. Accept solution (as Requester)
	fmt.Println("[6/6] Accepting solution (as Requester)...")
	_, err = registry.AcceptSolution(ctx, requestID, reward)
	if err != nil {
		t.Fatalf("Failed to accept solution: %v", err)
	}

	// Verify final status
	req, _ = registry.GetRequest(ctx, requestID)
	if req.Status != contracts.StatusAccepted {
		t.Fatalf("Expected status Accepted(2), got %d", req.Status)
	}
	fmt.Println("      Solution accepted!")

	fmt.Println()
	fmt.Println("=== E2E Test PASSED ===")
	fmt.Println()
	fmt.Println("Summary:")
	fmt.Printf("  - Request ID: %d\n", requestID)
	fmt.Printf("  - Reward: 10 SYRUP\n")
	fmt.Printf("  - Token Usage: %d tokens\n", tokenUsage)
	fmt.Printf("  - Final Status: Accepted\n")
}

func TestCancelRequest(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println()
	fmt.Println("=== E2E Test: Cancel Request ===")
	fmt.Println()

	client, err := ethclient.Dial(testRPC)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	registry, err := contracts.NewRegistryClient(testBakeRegistry, client, testPrivateKey)
	if err != nil {
		t.Fatalf("Failed to create registry client: %v", err)
	}

	// Create a request
	fmt.Println("[1/3] Creating request to cancel...")
	codeHash := contracts.HashCode("test code for cancel")
	reward := new(big.Int)
	reward.SetString("5000000000000000000", 10) // 5 SYRUP

	registry.ApproveToken(ctx, testSyrupToken, reward)
	_, requestID, err := registry.CreateRequest(ctx, codeHash, reward)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	fmt.Printf("      Request #%d created\n", requestID)

	// Cancel the request
	fmt.Println("[2/3] Cancelling request...")
	_, err = registry.CancelRequest(ctx, requestID)
	if err != nil {
		t.Fatalf("Failed to cancel: %v", err)
	}

	// Verify cancelled
	fmt.Println("[3/3] Verifying cancellation...")
	req, _ := registry.GetRequest(ctx, requestID)
	if req.Status != contracts.StatusCancelled {
		t.Fatalf("Expected status Cancelled(4), got %d", req.Status)
	}
	fmt.Println("      Request cancelled successfully!")
	fmt.Println()
	fmt.Println("=== Cancel Test PASSED ===")
}

func TestMain(m *testing.M) {
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("  WAFFLE END-TO-END TEST SUITE")
	fmt.Println("========================================")
	fmt.Println()
	os.Exit(m.Run())
}
