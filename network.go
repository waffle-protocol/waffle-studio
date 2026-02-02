package main

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/generative-ai-go/genai"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/sergi/go-diff/diffmatchpatch"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"
)

const (
	aesKey     = "thisis32bitlongpassphrasetouse!!"
	serviceTag = "syrup-p2p-service"
	msgTypeRequest        = "request"
	msgTypePaymentRequest = "payment_request"
	msgTypePaymentProof   = "payment_proof"
	msgTypePaymentFailed  = "payment_failed"
	msgTypeResult         = "result"
	msgTypeError          = "error"
	erc20BalanceOfABI     = `[{"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	erc20TransferABI      = `[{"inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]`
)

type SyrupPayload struct {
	Type            string `json:"type,omitempty"`
	RequestID       string `json:"request_id,omitempty"`
	Prompt          string `json:"prompt,omitempty"`
	Data            []byte `json:"data,omitempty"`
	TokenUsage      uint64 `json:"token_usage,omitempty"`
	PricingRate     uint64 `json:"pricing_rate,omitempty"`
	PriceWei        string `json:"price_wei,omitempty"`
	PriceSyrup      string `json:"price_syrup,omitempty"`
	ProviderAddress string `json:"provider_address,omitempty"`
	ClientAddress   string `json:"client_address,omitempty"`
	TxHash          string `json:"tx_hash,omitempty"`
	Error           string `json:"error,omitempty"`
}

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

type paymentConfig struct {
	RPCURL          string `yaml:"rpc_url"`
	PrivateKey      string `yaml:"private_key"`
	SyrupToken      string `yaml:"syrup_token"`
	PricingRate     uint64 `yaml:"pricing_rate"`
	ProviderAddress string `yaml:"provider_address"`
}

type syrupWallet struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address
	client     *ethclient.Client
	chainID    *big.Int
}

type pendingResult struct {
	encryptedResult []byte
	tokenUsage      uint64
	priceWei        *big.Int
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

func initMDNS(h host.Host, rendezvous string) (chan peer.AddrInfo, error) {
	n := &discoveryNotifee{PeerChan: make(chan peer.AddrInfo)}
	s := mdns.NewMdnsService(h, rendezvous, n)
	return n.PeerChan, s.Start()
}

func expandPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, path[2:]), nil
	}
	return path, nil
}

func loadPaymentConfig() (*paymentConfig, error) {
	cfg := &paymentConfig{
		PricingRate: 1,
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(homeDir, ".waffle", "config.yaml")
		if info, err := os.Stat(configPath); err == nil && !info.IsDir() {
			data, err := os.ReadFile(configPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("failed to parse config file: %w", err)
			}
		} else if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to stat config file: %w", err)
		}
	}

	if val := os.Getenv("RPC_URL"); val != "" {
		cfg.RPCURL = val
	}
	if val := os.Getenv("PRIVATE_KEY"); val != "" {
		cfg.PrivateKey = val
	}
	if val := os.Getenv("SYRUP_TOKEN"); val != "" {
		cfg.SyrupToken = val
	}
	if val := os.Getenv("PRICING_RATE"); val != "" {
		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid PRICING_RATE: %w", err)
		}
		cfg.PricingRate = parsed
	}
	if val := os.Getenv("PROVIDER_ADDRESS"); val != "" {
		cfg.ProviderAddress = val
	}

	if cfg.RPCURL == "" {
		cfg.RPCURL = "http://127.0.0.1:8545"
	}
	if cfg.PricingRate == 0 {
		cfg.PricingRate = 1
	}

	return cfg, nil
}

func deriveAddress(privateKeyHex string) (common.Address, error) {
	keyHex := strings.TrimPrefix(strings.TrimSpace(privateKeyHex), "0x")
	privateKey, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		return common.Address{}, fmt.Errorf("invalid private key: %w", err)
	}
	return crypto.PubkeyToAddress(privateKey.PublicKey), nil
}

func newSyrupWallet(cfg *paymentConfig) (*syrupWallet, error) {
	if cfg.PrivateKey == "" {
		return nil, fmt.Errorf("private key not set")
	}
	keyHex := strings.TrimPrefix(strings.TrimSpace(cfg.PrivateKey), "0x")
	privateKey, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &syrupWallet{
		privateKey: privateKey,
		address:    address,
		client:     client,
		chainID:    chainID,
	}, nil
}

func (w *syrupWallet) Close() {
	if w.client != nil {
		w.client.Close()
	}
}

func (w *syrupWallet) AddressHex() string {
	return w.address.Hex()
}

func (w *syrupWallet) TokenBalance(ctx context.Context, tokenAddress common.Address) (*big.Int, error) {
	parsedABI, err := abi.JSON(strings.NewReader(erc20BalanceOfABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := parsedABI.Pack("balanceOf", w.address)
	if err != nil {
		return nil, fmt.Errorf("failed to pack call data: %w", err)
	}

	result, err := w.client.CallContract(ctx, ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	var balance *big.Int
	if err := parsedABI.UnpackIntoInterface(&balance, "balanceOf", result); err != nil {
		return nil, fmt.Errorf("failed to unpack result: %w", err)
	}

	return balance, nil
}

func (w *syrupWallet) TransferToken(ctx context.Context, tokenAddress common.Address, to common.Address, amount *big.Int) (common.Hash, *types.Receipt, error) {
	parsedABI, err := abi.JSON(strings.NewReader(erc20TransferABI))
	if err != nil {
		return common.Hash{}, nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	data, err := parsedABI.Pack("transfer", to, amount)
	if err != nil {
		return common.Hash{}, nil, fmt.Errorf("failed to pack transfer data: %w", err)
	}

	nonce, err := w.client.PendingNonceAt(ctx, w.address)
	if err != nil {
		return common.Hash{}, nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := w.client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	gasLimit, err := w.client.EstimateGas(ctx, ethereum.CallMsg{
		From: w.address,
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		gasLimit = 120000
	}

	tx := types.NewTransaction(nonce, tokenAddress, big.NewInt(0), gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(w.chainID), w.privateKey)
	if err != nil {
		return common.Hash{}, nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	if err := w.client.SendTransaction(ctx, signedTx); err != nil {
		return common.Hash{}, nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	receipt, err := waitForReceipt(ctx, w.client, signedTx.Hash())
	if err != nil {
		return signedTx.Hash(), nil, err
	}

	return signedTx.Hash(), receipt, nil
}

func waitForReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
}

func parseBigInt(value string) (*big.Int, error) {
	if strings.TrimSpace(value) == "" {
		return nil, fmt.Errorf("empty amount")
	}
	amt, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount: %s", value)
	}
	return amt, nil
}

func formatSyrupAmount(wei *big.Int) string {
	if wei == nil {
		return "0"
	}
	r := new(big.Rat).SetInt(wei)
	r.Quo(r, new(big.Rat).SetInt(big.NewInt(1_000_000_000_000_000_000)))
	return r.FloatString(6)
}

func estimateTokenUsage(parts ...string) uint64 {
	totalChars := 0
	for _, part := range parts {
		totalChars += len(part)
	}
	estimated := totalChars / 4
	if estimated < 1 {
		estimated = 1
	}
	return uint64(estimated)
}

func newRequestID() string {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

func priceWeiFromUsage(tokenUsage uint64, rate uint64) *big.Int {
	usage := new(big.Int).SetUint64(tokenUsage)
	r := new(big.Int).SetUint64(rate)
	usage.Mul(usage, r)
	return usage.Mul(usage, big.NewInt(1_000_000_000_000_000_000))
}

func sendPayload(ctx context.Context, node host.Host, target peer.ID, payload SyrupPayload) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	s, err := node.NewStream(ctx, target, "/syrup/v1")
	if err != nil {
		return err
	}
	defer s.Close()
	_, err = s.Write(jsonData)
	return err
}

func verifyPaymentTx(ctx context.Context, client *ethclient.Client, tokenAddress common.Address, providerAddress common.Address, txHash common.Hash, expected *big.Int) (bool, error) {
	receipt, err := waitForReceipt(ctx, client, txHash)
	if err != nil {
		return false, err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return false, fmt.Errorf("transaction failed")
	}
	if expected == nil || expected.Sign() == 0 {
		return true, nil
	}

	transferSig := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	total := new(big.Int)

	for _, log := range receipt.Logs {
		if log.Address != tokenAddress {
			continue
		}
		if len(log.Topics) < 3 || log.Topics[0] != transferSig {
			continue
		}
		to := common.BytesToAddress(log.Topics[2].Bytes())
		if to != providerAddress {
			continue
		}
		if len(log.Data) == 0 {
			continue
		}
		amt := new(big.Int).SetBytes(log.Data)
		total.Add(total, amt)
	}

	if total.Cmp(expected) < 0 {
		return false, fmt.Errorf("transfer amount %s < expected %s", total.String(), expected.String())
	}

	return true, nil
}

func getAPIKey() string {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".syrup")
	keyFile := filepath.Join(configDir, "apikey")

	if data, err := os.ReadFile(keyFile); err == nil {
		return strings.TrimSpace(string(data))
	}

	fmt.Println("\n⚠️  Gemini API 키가 설정되지 않았습니다.")
	fmt.Print("🔑 API 키를 입력하세요 (화면에 표시됨): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	inputKey := strings.TrimSpace(scanner.Text())

	if inputKey == "" {
		fmt.Println("❌ 키 입력 취소. 종료합니다.")
		os.Exit(1)
	}
	os.MkdirAll(configDir, 0700)
	os.WriteFile(keyFile, []byte(inputKey), 0600)
	return inputKey
}

func encrypt(data []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	return gcm.Seal(nonce, nonce, data, nil)
}

func decrypt(data []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil
	}
	return plaintext
}

func printGitStyleStat(filename string, original, new string) {
	dmp := diffmatchpatch.New()
	text1, text2, lineArray := dmp.DiffLinesToChars(original, new)
	diffs := dmp.DiffMain(text1, text2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	insertions := 0
	deletions := 0

	for _, diff := range diffs {
		lines := strings.Count(diff.Text, "\n")
		if lines == 0 && len(diff.Text) > 0 {
			lines = 1
		}
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			insertions += lines
		case diffmatchpatch.DiffDelete:
			deletions += lines
		}
	}

	totalChanges := insertions + deletions
	if totalChanges == 0 {
		fmt.Println("No changes detected.")
		return
	}

	maxGraphWidth := 30
	graphWidth := totalChanges
	if graphWidth > maxGraphWidth {
		graphWidth = maxGraphWidth
	}

	plusCount := 0
	if totalChanges > 0 {
		plusCount = (insertions * graphWidth) / totalChanges
	}
	minusCount := graphWidth - plusCount

	displayFilename := filepath.Base(filename)
	if len(displayFilename) > 20 {
		displayFilename = displayFilename[:17] + "..."
	}

	fmt.Println("\n📝 [변경 사항 요약]")
	fmt.Printf(" %-20s | %d ", displayFilename, totalChanges)
	fmt.Printf("\033[32m%s\033[0m", strings.Repeat("+", plusCount))
	fmt.Printf("\033[31m%s\033[0m\n", strings.Repeat("-", minusCount))
	fmt.Printf(" 1 file changed, %d insertions(+), %d deletions(-)\n", insertions, deletions)
	fmt.Println()
}

func main() {
	mode := flag.String("mode", "client", "execution mode: client or provider")
	flag.Parse()

	ctx := context.Background()
	var model *genai.GenerativeModel
	responseChan := make(chan []byte)
	var payCfg *paymentConfig
	var wallet *syrupWallet
	var ethClient *ethclient.Client
	var syrupTokenAddr common.Address
	var providerAddress common.Address
	var pricingRate uint64
	pendingResults := make(map[string]pendingResult)
	var pendingMu sync.Mutex

	// ================= Provider 설정 =================
	if *mode == "provider" {
		apiKey := getAPIKey()
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			fmt.Println("❌ API 클라이언트 생성 실패:", err)
			return
		}
		model = client.GenerativeModel("models/gemini-2.5-flash")
		fmt.Println("🌟 Provider 모드: Gemini 엔진 준비 완료.")

		payCfg, err = loadPaymentConfig()
		if err != nil {
			fmt.Println("❌ 결제 설정 로드 실패:", err)
			return
		}
		if payCfg.SyrupToken == "" || !common.IsHexAddress(payCfg.SyrupToken) {
			fmt.Println("❌ SYRUP_TOKEN 설정이 필요합니다.")
			return
		}
		syrupTokenAddr = common.HexToAddress(payCfg.SyrupToken)
		pricingRate = payCfg.PricingRate

		if payCfg.ProviderAddress != "" {
			if !common.IsHexAddress(payCfg.ProviderAddress) {
				fmt.Println("❌ PROVIDER_ADDRESS 형식이 올바르지 않습니다.")
				return
			}
			providerAddress = common.HexToAddress(payCfg.ProviderAddress)
		} else {
			if payCfg.PrivateKey == "" {
				fmt.Println("❌ PRIVATE_KEY 설정이 필요합니다. (provider)")
				return
			}
			providerAddr, err := deriveAddress(payCfg.PrivateKey)
			if err != nil {
				fmt.Println("❌ Provider 주소 생성 실패:", err)
				return
			}
			providerAddress = providerAddr
		}

		ethClient, err = ethclient.Dial(payCfg.RPCURL)
		if err != nil {
			fmt.Println("❌ RPC 연결 실패:", err)
			return
		}
		defer ethClient.Close()

		fmt.Printf("💳 결제 활성화: Provider %s | Rate %d SYRUP/token\n", providerAddress.Hex(), pricingRate)
	}

	// ================= Client 결제 설정 =================
	if *mode == "client" {
		var err error
		payCfg, err = loadPaymentConfig()
		if err != nil {
			fmt.Println("❌ 결제 설정 로드 실패:", err)
			return
		}
		if payCfg.SyrupToken == "" || !common.IsHexAddress(payCfg.SyrupToken) {
			fmt.Println("❌ SYRUP_TOKEN 설정이 필요합니다.")
			return
		}
		syrupTokenAddr = common.HexToAddress(payCfg.SyrupToken)
		if payCfg.PrivateKey == "" {
			fmt.Println("❌ PRIVATE_KEY 설정이 필요합니다. (client)")
			return
		}
		wallet, err = newSyrupWallet(payCfg)
		if err != nil {
			fmt.Println("❌ 지갑 초기화 실패:", err)
			return
		}
		defer wallet.Close()
	}

	// ================= P2P 노드 생성 =================
	node, _ := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	defer node.Close()

	peerChan, err := initMDNS(node, serviceTag)
	if err != nil {
		panic(err)
	}

	// ================= 데이터 핸들러 =================
	node.SetStreamHandler("/syrup/v1", func(s network.Stream) {
		defer s.Close()
		raw, _ := io.ReadAll(s)

		if *mode == "client" {
			responseChan <- raw
			return
		}

		if *mode != "provider" || model == nil {
			return
		}

		var payload SyrupPayload
		if err := json.Unmarshal(raw, &payload); err != nil {
			return
		}
		if payload.Type == "" {
			payload.Type = msgTypeRequest
		}

		switch payload.Type {
		case msgTypeRequest:
			if payload.RequestID == "" {
				payload.RequestID = newRequestID()
			}

			decryptedInput := decrypt(payload.Data, []byte(aesKey))
			if decryptedInput == nil {
				_ = sendPayload(ctx, node, s.Conn().RemotePeer(), SyrupPayload{
					Type:      msgTypeError,
					RequestID: payload.RequestID,
					Error:     "decrypt failed",
				})
				return
			}

			fmt.Printf("\n⚙️  요청 수신: \"%s\"\n", payload.Prompt)
			prompt := fmt.Sprintf("User Request: %s\n\nRewrite the following file content to meet the request. Output ONLY the raw content code/text. Do NOT use markdown code blocks (```) or explanation:\n\n%s", payload.Prompt, string(decryptedInput))

			resp, err := model.GenerateContent(ctx, genai.Text(prompt))
			if err != nil {
				_ = sendPayload(ctx, node, s.Conn().RemotePeer(), SyrupPayload{
					Type:      msgTypeError,
					RequestID: payload.RequestID,
					Error:     fmt.Sprintf("AI 처리 에러: %v", err),
				})
				fmt.Println("❌ AI 처리 에러:", err)
				return
			}

			var aiResult string
			for _, cand := range resp.Candidates {
				for _, part := range cand.Content.Parts {
					aiResult += fmt.Sprintf("%v", part)
				}
			}
			aiResult = strings.TrimPrefix(aiResult, "```go")
			aiResult = strings.TrimPrefix(aiResult, "```")
			aiResult = strings.TrimSuffix(aiResult, "```")

			tokenUsage := uint64(0)
			if resp.UsageMetadata != nil {
				tokenUsage = uint64(resp.UsageMetadata.TotalTokenCount)
			}
			if tokenUsage == 0 {
				tokenUsage = estimateTokenUsage(string(decryptedInput), aiResult)
				fmt.Printf("⚠️ 토큰 사용량 추정값 사용: %d\n", tokenUsage)
			}

			priceWei := priceWeiFromUsage(tokenUsage, pricingRate)
			encryptedResponse := encrypt([]byte(aiResult), []byte(aesKey))

			if priceWei.Sign() == 0 {
				_ = sendPayload(ctx, node, s.Conn().RemotePeer(), SyrupPayload{
					Type:       msgTypeResult,
					RequestID:  payload.RequestID,
					Data:       encryptedResponse,
					TokenUsage: tokenUsage,
				})
				fmt.Println("✅ 처리 완료 및 전송됨! (무료)")
				return
			}

			pendingMu.Lock()
			pendingResults[payload.RequestID] = pendingResult{
				encryptedResult: encryptedResponse,
				tokenUsage:      tokenUsage,
				priceWei:        priceWei,
			}
			pendingMu.Unlock()

			paymentPayload := SyrupPayload{
				Type:            msgTypePaymentRequest,
				RequestID:       payload.RequestID,
				TokenUsage:      tokenUsage,
				PricingRate:     pricingRate,
				PriceWei:        priceWei.String(),
				PriceSyrup:      formatSyrupAmount(priceWei),
				ProviderAddress: providerAddress.Hex(),
			}

			if err := sendPayload(ctx, node, s.Conn().RemotePeer(), paymentPayload); err != nil {
				fmt.Println("❌ 결제 요청 전송 실패:", err)
				return
			}
			fmt.Println("💳 결제 요청 전송됨.")

		case msgTypePaymentProof:
			if payload.RequestID == "" {
				return
			}

			pendingMu.Lock()
			pending, ok := pendingResults[payload.RequestID]
			pendingMu.Unlock()
			if !ok {
				return
			}
			if payload.TxHash == "" {
				_ = sendPayload(ctx, node, s.Conn().RemotePeer(), SyrupPayload{
					Type:      msgTypeError,
					RequestID: payload.RequestID,
					Error:     "tx hash missing",
				})
				return
			}

			verifyCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
			defer cancel()

			okPay, err := verifyPaymentTx(verifyCtx, ethClient, syrupTokenAddr, providerAddress, common.HexToHash(payload.TxHash), pending.priceWei)
			if err != nil || !okPay {
				_ = sendPayload(ctx, node, s.Conn().RemotePeer(), SyrupPayload{
					Type:      msgTypeError,
					RequestID: payload.RequestID,
					Error:     fmt.Sprintf("결제 확인 실패: %v", err),
				})
				fmt.Println("❌ 결제 확인 실패:", err)
				return
			}

			if err := sendPayload(ctx, node, s.Conn().RemotePeer(), SyrupPayload{
				Type:       msgTypeResult,
				RequestID:  payload.RequestID,
				Data:       pending.encryptedResult,
				TokenUsage: pending.tokenUsage,
			}); err != nil {
				fmt.Println("❌ 결과 전송 실패:", err)
				return
			}

			pendingMu.Lock()
			delete(pendingResults, payload.RequestID)
			pendingMu.Unlock()

			fmt.Println("✅ 결제 확인 및 코드 전송됨!")

		case msgTypePaymentFailed:
			if payload.RequestID != "" {
				pendingMu.Lock()
				delete(pendingResults, payload.RequestID)
				pendingMu.Unlock()
			}
			if payload.Error != "" {
				fmt.Printf("❌ 결제 실패: %s\n", payload.Error)
			}
		}
	})

	fmt.Printf("\n🚀 [%s] 실행 중 (ID: %s)\n", strings.ToUpper(*mode), node.ID())

	// ================= [Provider] 대기 모드 =================
	if *mode == "provider" {
		fmt.Println("📡 네트워크에 내 존재를 알리는 중... (Client 검색 대기)")
		select {}
	}

	// ================= [Client] 검색 로직 (수정됨) =================
	var connectedPeer peer.ID
	var foundPeers []peer.AddrInfo

	fmt.Println("🔍 주변 Provider 검색 중... (Provider가 켜질 때까지 대기합니다)")

	// 1단계: 첫 번째 Provider가 나올 때까지 무한 대기 (Blocking)
	for {
		peerInfo := <-peerChan
		if peerInfo.ID != node.ID() {
			foundPeers = append(foundPeers, peerInfo)
			fmt.Printf("   Found: %s\n", peerInfo.ID)
			break // 한 명 찾았으니 루프 탈출하고 추가 수집으로 이동
		}
	}

	// 2단계: 1.5초 더 기다려서 추가 Provider 수집 (동시에 켜져 있을 경우 대비)
	fmt.Println("   (추가 Provider 확인 중...)")
	timeout := time.After(1500 * time.Millisecond)

CollectionLoop:
	for {
		select {
		case peerInfo := <-peerChan:
			if peerInfo.ID == node.ID() {
				continue
			}

			// 중복 체크
			isNew := true
			for _, p := range foundPeers {
				if p.ID == peerInfo.ID {
					isNew = false
					break
				}
			}
			if isNew {
				foundPeers = append(foundPeers, peerInfo)
				fmt.Printf("   Found: %s\n", peerInfo.ID)
			}
		case <-timeout:
			break CollectionLoop
		}
	}

	// 목록 출력 및 선택
	fmt.Println("\n📋 [발견된 Provider 목록]")
	for i, p := range foundPeers {
		fmt.Printf("[%d] %s\n", i+1, p.ID)
	}

	fmt.Print("\n번호를 선택하세요: ")
	var selection int
	_, err = fmt.Scanln(&selection)
	if err != nil || selection < 1 || selection > len(foundPeers) {
		fmt.Println("❌ 잘못된 입력입니다. 종료합니다.")
		return
	}

	targetPeer := foundPeers[selection-1]
	fmt.Printf("🔗 Provider [%s]에 연결 시도 중...\n", targetPeer.ID)

	if err := node.Connect(ctx, targetPeer); err != nil {
		fmt.Println("❌ 연결 실패:", err)
		return
	}

	connectedPeer = targetPeer.ID
	fmt.Println("✅ 연결 성공!")

	// 버퍼 비우기용 스캐너 재설정
	scanner := bufio.NewScanner(os.Stdin)

	// ================= [Client] 파일 전송 루프 =================
	for {
		fmt.Print("\n📂 수정할 파일 경로 (예: ~/Desktop/a.txt): ")
		if !scanner.Scan() {
			break
		}
		rawPath := scanner.Text()

		if strings.TrimSpace(rawPath) == "" {
			continue
		}

		targetFile, err := expandPath(rawPath)
		if err != nil {
			fmt.Println("❌ 경로 에러:", err)
			continue
		}

		fileData, err := os.ReadFile(targetFile)
		if err != nil {
			fmt.Printf("❌ '%s' 파일을 찾을 수 없습니다.\n", targetFile)
			continue
		}

		fmt.Print("📝 요청사항 입력: ")
		scanner.Scan()
		prompt := scanner.Text()

		encryptedData := encrypt(fileData, []byte(aesKey))
		requestID := newRequestID()
		payload := SyrupPayload{
			Type:      msgTypeRequest,
			RequestID: requestID,
			Prompt:    prompt,
			Data:      encryptedData,
		}

		if err := sendPayload(ctx, node, connectedPeer, payload); err != nil {
			fmt.Println("❌ 전송 실패 (연결 끊김?):", err)
			break
		}
		fmt.Println("🚀 전송 완료! 결제/응답 대기 중...")

		var decryptedResult []byte
		gotResult := false
		abort := false

		for {
			respRaw := <-responseChan
			var respPayload SyrupPayload
			if err := json.Unmarshal(respRaw, &respPayload); err != nil {
				continue
			}
			if respPayload.RequestID != "" && respPayload.RequestID != requestID {
				continue
			}
			if respPayload.Type == "" {
				respPayload.Type = msgTypeResult
			}

			failPayment := func(msg string) {
				_ = sendPayload(ctx, node, connectedPeer, SyrupPayload{
					Type:      msgTypePaymentFailed,
					RequestID: requestID,
					Error:     msg,
				})
				fmt.Println("❌ 결제 실패:", msg)
				abort = true
			}

			switch respPayload.Type {
			case msgTypePaymentRequest:
				if respPayload.PriceWei == "" {
					failPayment("결제 금액 누락")
					break
				}
				priceWei, err := parseBigInt(respPayload.PriceWei)
				if err != nil {
					failPayment("결제 금액 파싱 실패")
					break
				}
				if respPayload.ProviderAddress == "" || !common.IsHexAddress(respPayload.ProviderAddress) {
					failPayment("Provider 주소 오류")
					break
				}
				providerAddr := common.HexToAddress(respPayload.ProviderAddress)

				balance, err := wallet.TokenBalance(ctx, syrupTokenAddr)
				if err != nil {
					failPayment("잔액 조회 실패")
					break
				}
				if balance.Cmp(priceWei) < 0 {
					failPayment("잔액 부족")
					break
				}

				fmt.Printf("💳 결제 진행: %s SYRUP (사용량: %d, Rate: %d)\n", respPayload.PriceSyrup, respPayload.TokenUsage, respPayload.PricingRate)
				payCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
				txHash, receipt, err := wallet.TransferToken(payCtx, syrupTokenAddr, providerAddr, priceWei)
				cancel()
				if err != nil {
					failPayment("결제 트랜잭션 실패")
					break
				}
				if receipt == nil || receipt.Status != types.ReceiptStatusSuccessful {
					failPayment("결제 트랜잭션 실패")
					break
				}

				proofPayload := SyrupPayload{
					Type:          msgTypePaymentProof,
					RequestID:     requestID,
					TxHash:        txHash.Hex(),
					PriceWei:      priceWei.String(),
					ClientAddress: wallet.AddressHex(),
				}
				if err := sendPayload(ctx, node, connectedPeer, proofPayload); err != nil {
					failPayment("결제 증빙 전송 실패")
					break
				}
				fmt.Println("✅ 결제 완료! 코드 수신 대기 중...")

			case msgTypeResult:
				decryptedResult = decrypt(respPayload.Data, []byte(aesKey))
				if decryptedResult == nil {
					fmt.Println("❌ 응답 복호화 실패")
					abort = true
					break
				}
				gotResult = true

			case msgTypeError:
				if respPayload.Error != "" {
					fmt.Printf("❌ 오류: %s\n", respPayload.Error)
				} else {
					fmt.Println("❌ 오류 발생")
				}
				abort = true

			case msgTypePaymentFailed:
				if respPayload.Error != "" {
					fmt.Printf("❌ 결제 실패: %s\n", respPayload.Error)
				} else {
					fmt.Println("❌ 결제 실패")
				}
				abort = true
			}

			if gotResult || abort {
				break
			}
		}

		if abort {
			continue
		}

		newContent := string(decryptedResult)
		originalContent := string(fileData)

		printGitStyleStat(targetFile, originalContent, newContent)

		fmt.Print("💡 변경사항을 적용하시겠습니까? (y/n): ")
		scanner.Scan()
		if strings.ToLower(strings.TrimSpace(scanner.Text())) == "y" {
			os.WriteFile(targetFile, decryptedResult, 0644)
			fmt.Printf("✨ [완료] 업데이트됨!\n")
		} else {
			fmt.Println("❌ [취소됨]")
		}
	}
}
