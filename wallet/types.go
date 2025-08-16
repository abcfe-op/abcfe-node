package wallet

import (
	prt "github.com/abcfe/abcfe-node/protocol"
)

// 기존 키스토어 관련 타입들
type CipherParams struct {
	IV string `json:"iv"` // 초기화 벡터
}

type KDFParams struct {
	DkLen int    `json:"dklen"` // 도출된 키 길이
	N     int    `json:"n"`     // CPU/메모리 비용
	P     int    `json:"p"`     // 병렬화 매개변수
	R     int    `json:"r"`     // 블록 크기
	Salt  string `json:"salt"`  // 솔트
}

type Crypto struct {
	Cipher       string       `json:"cipher"`     // "aes-128-ctr"
	CipherText   string       `json:"ciphertext"` // 암호화된 개인키
	CipherParams CipherParams `json:"cipherparams"`
	KDF          string       `json:"kdf"` // "scrypt"
	KDFParams    KDFParams    `json:"kdfparams"`
	MAC          string       `json:"mac"` // 무결성 검증
}

// 니모닉 기반 지갑 타입들
type MnemonicWallet struct {
	Mnemonic     string     `json:"mnemonic"`      // 12/15/18/21/24개 단어
	Seed         []byte     `json:"seed"`          // 니모닉에서 도출된 시드
	MasterKey    []byte     `json:"master_key"`    // 마스터 개인키 (바이트)
	Accounts     []*Account `json:"accounts"`      // 파생된 계정들
	CurrentIndex int        `json:"current_index"` // 현재 사용 중인 계정 인덱스
}

type Account struct {
	Index      int         `json:"index"`       // 계정 인덱스 (0, 1, 2...)
	Address    prt.Address `json:"address"`     // 20바이트 주소
	PrivateKey []byte      `json:"private_key"` // 개인키 (바이트)
	PublicKey  []byte      `json:"public_key"`  // 공개키 (바이트)
	Path       string      `json:"path"`        // BIP-44 경로 (m/44'/60'/0'/0/0)
	Unlocked   bool        `json:"unlocked"`    // 언락 상태
}

// BIP-44 경로 상수
const (
	BIP44Purpose  = 44
	BIP44CoinType = 60 // Ethereum
	BIP44Account  = 0
	BIP44Change   = 0 // External
	BIP44Index    = 0
)
