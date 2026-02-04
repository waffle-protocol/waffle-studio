package ui

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

// BoxWidth는 기본 너비입니다.
const BoxWidth = 64

type Box struct {
	width int
}

func NewBox() *Box {
	return &Box{width: BoxWidth}
}

func NewBoxWithWidth(width int) *Box {
	return &Box{width: width}
}

// Header: 상단 테두리와 제목 (정렬 보정)
func (b *Box) Header(title string) {
	top := "╔" + strings.Repeat("═", b.width-2) + "╗"
	fmt.Printf("%s%s%s\n", Blue, top, Reset)

	// 제목도 Row와 동일한 간격 로직 적용
	content := Amber + Bold + title + Reset
	b.Row(content)

	mid := "╠" + strings.Repeat("═", b.width-2) + "╣"
	fmt.Printf("%s%s%s\n", Blue, mid, Reset)
}

// Row: 모든 행의 우측 세로선 위치를 동일하게 맞춤
func (b *Box) Row(content string) {
	visibleLen := displayWidth(content)
	// ║ (1) + 공백 (1) + 내용 + 공백 (n) + ║ (1) 구조
	// b.width에서 세로선 2개와 좌측 공백 1개를 뺀 나머지 공간 계산
	padding := b.width - visibleLen - 3
	if padding < 0 {
		padding = 0
	}

	fmt.Printf("%s║ %s%s%s%s║%s\n", Blue, Reset, content, strings.Repeat(" ", padding), Blue, Reset)
}

func (b *Box) Footer() {
	bottom := "╚" + strings.Repeat("═", b.width-2) + "╝"
	fmt.Printf("%s%s%s\n", Blue, bottom, Reset)
}

// WalletBox: 출력 데이터 구성
func WalletBox(address string, balance float64, tokenName string) {
	addressRow := fmt.Sprintf("💳 Address: %s%s%s", Amber, address, Reset)
	balanceRow := fmt.Sprintf("🍯 Balance: %s%.4f%s", Amber, balance, Reset)
	tokenRow := fmt.Sprintf("            %s%s%s", Bold, tokenName, Reset)

	// 가장 긴 줄에 맞춰 박스 너비 유동적 조절
	width := BoxWidth
	for _, row := range []string{addressRow, balanceRow, tokenRow, "WAFFLE WALLET"} {
		if displayWidth(row)+6 > width {
			width = displayWidth(row) + 6
		}
	}

	box := NewBoxWithWidth(width)
	fmt.Println()
	box.Header("WAFFLE WALLET")
	box.Row(addressRow)
	box.Row(balanceRow)
	box.Row(tokenRow)
	box.Footer()
	fmt.Println()
}

// displayWidth: ANSI 제외 실제 너비 계산
func displayWidth(s string) int {
	res := 0
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		res += runewidth.RuneWidth(r)
	}
	return res
}

// 나머지 헬퍼 함수들
func (b *Box) Separator() {
	fmt.Printf("%s╠%s╣%s\n", Blue, strings.Repeat("═", b.width-2), Reset)
}

func ResultBox(title, requestID string, fields map[string]string, txHash string) {
	box := NewBox()
	fmt.Println()
	box.Header(title)
	if requestID != "" {
		box.Row(fmt.Sprintf("Request ID: #%s", requestID))
	}
	for l, v := range fields {
		box.Row(fmt.Sprintf("%s: %s%s%s", l, Amber, v, Reset))
	}
	if txHash != "" {
		box.Row(fmt.Sprintf("🔗 Tx: %s", txHash))
	}
	box.Footer()
	fmt.Println()
}

func RequestsHeader(title string, width int) {
	box := NewBoxWithWidth(width)
	fmt.Println()
	box.Header(title)
}
func RequestsFooter(width int)              { box := NewBoxWithWidth(width); box.Footer(); fmt.Println() }
func RequestsRow(width int, content string) { box := NewBoxWithWidth(width); box.Row(content) }
func truncateHash(hash string, maxLen int) string {
	if len(hash) <= maxLen {
		return hash
	}
	return hash[:maxLen]
}
