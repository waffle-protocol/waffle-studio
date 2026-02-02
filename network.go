package main

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/sergi/go-diff/diffmatchpatch"
	"google.golang.org/api/option"
)

const (
	aesKey     = "thisis32bitlongpassphrasetouse!!"
	serviceTag = "syrup-p2p-service"
)

type SyrupPayload struct {
	Prompt string `json:"prompt,omitempty"`
	Data   []byte `json:"data"`
}

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
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

		if *mode == "provider" && model != nil {
			var payload SyrupPayload
			if err := json.Unmarshal(raw, &payload); err != nil {
				return
			}

			decryptedInput := decrypt(payload.Data, []byte(aesKey))
			fmt.Printf("\n⚙️  요청 수신: \"%s\"\n", payload.Prompt)

			prompt := fmt.Sprintf("User Request: %s\n\nRewrite the following file content to meet the request. Output ONLY the raw content code/text. Do NOT use markdown code blocks (```) or explanation:\n\n%s", payload.Prompt, string(decryptedInput))

			resp, err := model.GenerateContent(ctx, genai.Text(prompt))
			var aiResult string
			if err == nil {
				for _, cand := range resp.Candidates {
					for _, part := range cand.Content.Parts {
						aiResult += fmt.Sprintf("%v", part)
					}
				}
				aiResult = strings.TrimPrefix(aiResult, "```go")
				aiResult = strings.TrimPrefix(aiResult, "```")
				aiResult = strings.TrimSuffix(aiResult, "```")
			} else {
				aiResult = string(decryptedInput)
				fmt.Println("❌ AI 처리 에러:", err)
			}

			encryptedResponse := encrypt([]byte(aiResult), []byte(aesKey))
			respPayload := SyrupPayload{Data: encryptedResponse}
			respJson, _ := json.Marshal(respPayload)

			sOut, _ := node.NewStream(ctx, s.Conn().RemotePeer(), "/syrup/v1")
			sOut.Write(respJson)
			sOut.Close()
			fmt.Println("✅ 처리 완료 및 전송됨!")
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
		payload := SyrupPayload{Prompt: prompt, Data: encryptedData}
		jsonData, _ := json.Marshal(payload)

		s, err := node.NewStream(ctx, connectedPeer, "/syrup/v1")
		if err != nil {
			fmt.Println("❌ 전송 실패 (연결 끊김?):", err)
			break
		}
		s.Write(jsonData)
		s.Close()
		fmt.Println("🚀 전송 완료! AI 응답 대기 중...")

		respRaw := <-responseChan

		var respPayload SyrupPayload
		json.Unmarshal(respRaw, &respPayload)
		decryptedResult := decrypt(respPayload.Data, []byte(aesKey))
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
