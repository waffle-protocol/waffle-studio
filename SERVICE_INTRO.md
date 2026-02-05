# Waffle - 서비스 소개

## 한 줄 소개
**Waffle**은 P2P 기반의 탈중앙화 AI 코드 수정 마켓플레이스입니다.

---

## 서비스 개요

### 무엇을 하는 서비스인가요?
Waffle은 사용자가 코드 수정을 요청하면, "Baker"라 불리는 프로바이더가 AI를 활용해 코드를 수정하고, 그 대가로 **SYRUP 토큰**을 보상받는 P2P 마켓플레이스입니다.

### 핵심 가치
| 가치 | 설명 |
|------|------|
| **탈중앙화** | 중앙 서버 없이 P2P 네트워크로 직접 연결 |
| **신뢰성** | 블록체인 에스크로우로 안전한 거래 보장 |
| **접근성** | CLI 도구로 어디서든 간편하게 사용 |

---

## 작동 방식

```
┌──────────────┐         P2P          ┌──────────────┐
│   Requester  │ ◄──────────────────► │    Baker     │
│   (사용자)    │                      │  (프로바이더) │
└──────┬───────┘                      └──────┬───────┘
       │                                     │
       │  SYRUP 에스크로우                    │  AI 코드 수정
       ▼                                     ▼
┌─────────────────────────────────────────────────────┐
│              Blockchain (Base Sepolia)              │
│         BakeRegistry + SYRUP Token (ERC-20)         │
└─────────────────────────────────────────────────────┘
```

### 사용 흐름
1. **요청 생성** - 사용자가 `waffle bake`로 파일 선택 후 수정 요청
2. **에스크로우** - SYRUP 토큰이 스마트 컨트랙트에 예치
3. **AI 처리** - Baker가 Gemini API를 통해 코드 수정
4. **검토 및 결제** - 사용자가 결과 확인 후 수락/거부 결정

---

## 주요 기능

### CLI 명령어
| 명령어 | 설명 |
|--------|------|
| `waffle bake` | 코드 수정 요청 생성 |
| `waffle serve` | Baker 노드 실행 |
| `waffle accept` | 솔루션 수락 및 결제 |
| `waffle reject` | 솔루션 거부 |
| `waffle cancel` | 요청 취소 및 환급 |
| `waffle balance` | SYRUP 잔액 조회 |
| `waffle faucet` | 테스트넷 토큰 받기 |

---

## 기술 스택

### Backend
- **Language**: Go 1.25
- **P2P**: libp2p (mDNS 피어 발견)
- **AI**: Google Gemini API (gemini-2.5-flash)

### Blockchain
- **Network**: Base Sepolia Testnet
- **Smart Contract**: Solidity ^0.8.20
- **Token**: ERC-20 (SYRUP)
- **Security**: OpenZeppelin ReentrancyGuard

### Interface
- **CLI Framework**: Cobra
- **TUI**: Charmbracelet (Bubble Tea)

---

## 토큰 이코노미

### SYRUP Token
| 항목 | 내용 |
|------|------|
| 심볼 | SYRUP |
| 표준 | ERC-20 |
| 초기 공급량 | 1,000,000 SYRUP |
| 용도 | AI 코드 수정 서비스 결제 |

### 결제 흐름
```
사용자 → 에스크로우 예치 → Baker 작업 완료 → 사용자 승인 → Baker 지급
                                            ↓
                                      거부/취소 시 환급
```

---

## 아키텍처

### 프로젝트 구조
```
waffle/
├── cmd/waffle/        # CLI 커맨드
├── internal/
│   ├── p2p/           # P2P 네트워킹
│   ├── ai/            # Gemini AI 통합
│   ├── contracts/     # 블록체인 연동
│   ├── wallet/        # 지갑 관리
│   └── ui/            # TUI 컴포넌트
├── contracts/         # Solidity 스마트 계약
└── relay/             # 가스리스 릴레이 서버
```

### 스마트 컨트랙트
| 컨트랙트 | 역할 |
|----------|------|
| **BakeRegistry** | 요청/솔루션 관리, 에스크로우 |
| **SyrupToken** | ERC-20 토큰, Faucet |

---

## 차별점

| 기존 서비스 | Waffle |
|-------------|--------|
| 중앙 서버 의존 | P2P 탈중앙화 |
| 신용 기반 결제 | 블록체인 에스크로우 |
| 폐쇄적 API | 누구나 Baker가 될 수 있음 |
| 고정 가격 | 시장 기반 동적 가격 |

---

## 배포 정보

| 항목 | 값 |
|------|-----|
| 네트워크 | Base Sepolia Testnet |
| SYRUP Token | `0xb84284ddab9f7e2b14ba81c6f44db8d99488e23b` |
| BakeRegistry | `0x8c03552b7ae490ddc2e6b5a1b3452129e1135323` |
| RPC | https://sepolia.base.org |

---

## 시작하기

```bash
# 1. 설치
go install github.com/waffle-studio/waffle@latest

# 2. 테스트넷 토큰 받기
waffle faucet

# 3. 잔액 확인
waffle balance

# 4. 코드 수정 요청
waffle bake
```

---

## 라이선스
MIT License
