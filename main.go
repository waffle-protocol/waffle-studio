package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// AES-GCM 암호화 함수
func encrypt(data []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	return gcm.Seal(nonce, nonce, data, nil)
}

func main() {
	// 1. 파일 읽기 (a.txt와 b.txt가 미리 있어야 합니다)
	fileA, err := os.ReadFile("a.txt")
	fileB, err := os.ReadFile("b.txt")
	if err != nil {
		fmt.Println("❌ 파일 읽기 실패: a.txt와 b.txt를 먼저 만드세요.")
		return
	}

	// 2. Diff 및 Patch 생성
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(fileA), string(fileB), false)
	patches := dmp.PatchMake(string(fileA), diffs)
	patchText := dmp.PatchToText(patches)

	// 3. 암호화 (32바이트 키)
	key := []byte("thisis32bitlongpassphrasetouse!!")
	encryptedData := encrypt([]byte(patchText), key)

	// 4. 암호화된 파일 저장
	err = os.WriteFile("diff.patch.encrypted", encryptedData, 0644)
	if err != nil {
		fmt.Println("❌ 파일 저장 실패:", err)
		return
	}

	fmt.Println("✅ diff.patch.encrypted 파일이 성공적으로 생성되었습니다!")
}
