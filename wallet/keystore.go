package wallet

import (
	"fmt"
)

type KeyStore struct {
	Address string `json:"address"`
	Crypto  Crypto `json:"crypto"`  // types.go에서 정의된 Crypto 사용
	ID      string `json:"id"`      // UUID
	Version int    `json:"version"` // 키스토어 버전
}

// ECDSA 개인키 → AES-128-CTR 암호화 → 키스토어 저장
// 패스워드 → scrypt → 키 도출 → AES-128-CTR로 개인키 암호화

type KeyStoreManager struct {
	keystoreDir string
}

func NewKeyStoreManager(keystoreDir string) *KeyStoreManager {
	return &KeyStoreManager{
		keystoreDir: keystoreDir,
	}
}

func (ks *KeyStoreManager) NewAccount(password string) (*Account, error) {
	// TODO: ECDSA 키 생성 및 키스토어 저장 구현
	return nil, fmt.Errorf("not implemented yet")
}

func (ks *KeyStoreManager) LoadAccount(path, password string) (*Account, error) {
	// TODO: 키스토어 파일에서 계정 로드 구현
	return nil, fmt.Errorf("not implemented yet")
}

func (ks *KeyStoreManager) SaveAccount(account *Account, path string) error {
	// TODO: 계정을 키스토어 파일로 저장 구현
	return fmt.Errorf("not implemented yet")
}

func (ks *KeyStoreManager) Unlock(password string) error {
	// TODO: 계정 언락 구현
	return fmt.Errorf("not implemented yet")
}

func (ks *KeyStoreManager) Lock() {
	// TODO: 계정 락 구현
}
