package wallet

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/abcfe/abcfe-node/common/crypto"
	"github.com/tyler-smith/go-bip39"
)

// 더 명확한 인터페이스명
type WalletStorage interface {
	Exists(path string) bool
	Write(path string, data []byte) error
	Read(path string) ([]byte, error)
	CreateDir(path string) error
}

// 파일 시스템 구현
type FileWalletStorage struct{}

func (f FileWalletStorage) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (f FileWalletStorage) Write(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func (f FileWalletStorage) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (f FileWalletStorage) CreateDir(path string) error {
	return os.MkdirAll(path, 0700)
}

// WalletManager에 주입
type WalletManager struct {
	walletDir string
	wallet    *MnemonicWallet
	storage   WalletStorage
}

func NewWalletManager(walletDir string) *WalletManager {
	return &WalletManager{
		walletDir: walletDir,
		storage:   FileWalletStorage{},
	}
}

// 새 니모닉 지갑 생성
func (p *WalletManager) CreateWallet() (*MnemonicWallet, error) {
	// 128비트 엔트로피로 12개 단어 생성
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entropy: %w", err)
	}

	// 엔트로피에서 니모닉 생성
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	// 니모닉에서 시드 생성
	seed := bip39.NewSeed(mnemonic, "")

	// 시드에서 마스터 키 생성
	masterKey, err := crypto.DeriveMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to derive master key: %w", err)
	}

	wallet := &MnemonicWallet{
		Mnemonic:     mnemonic,
		Seed:         seed,
		MasterKey:    masterKey,
		Accounts:     []*Account{},
		CurrentIndex: 0,
	}

	// 첫 번째 계정 생성
	account, err := p.deriveAccount(wallet, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive first account: %w", err)
	}

	wallet.Accounts = append(wallet.Accounts, account)
	p.wallet = wallet

	return wallet, nil
}

// 기존 니모닉으로 지갑 복구
func (p *WalletManager) RestoreWallet(mnemonic string) (*MnemonicWallet, error) {
	// 니모닉 유효성 검사
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	// 니모닉에서 시드 생성
	seed := bip39.NewSeed(mnemonic, "")

	// 시드에서 마스터 키 생성
	masterKey, err := crypto.DeriveMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to derive master key: %w", err)
	}

	wallet := &MnemonicWallet{
		Mnemonic:     mnemonic,
		Seed:         seed,
		MasterKey:    masterKey,
		Accounts:     []*Account{},
		CurrentIndex: 0,
	}

	// 첫 번째 계정 생성
	account, err := p.deriveAccount(wallet, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive first account: %w", err)
	}

	wallet.Accounts = append(wallet.Accounts, account)
	p.wallet = wallet

	return wallet, nil
}

// 새 계정 추가
func (p *WalletManager) AddAccount() (*Account, error) {
	if p.wallet == nil {
		return nil, fmt.Errorf("wallet not initialized")
	}

	nextIndex := len(p.wallet.Accounts)
	account, err := p.deriveAccount(p.wallet, nextIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account: %w", err)
	}

	p.wallet.Accounts = append(p.wallet.Accounts, account)
	return account, nil
}

// 계정 파생 (BIP-44 경로: m/44'/60'/0'/0/index)
func (p *WalletManager) deriveAccount(wallet *MnemonicWallet, index int) (*Account, error) {
	path := fmt.Sprintf("m/%d'/%d'/%d'/%d/%d",
		BIP44Purpose, BIP44CoinType, BIP44Account, BIP44Change, index)

	// 마스터 키에서 계정 키 파생
	privateKey, publicKey, err := crypto.DeriveAccountKey(wallet.MasterKey, path)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account key: %w", err)
	}

	// 공개키에서 주소 생성
	address, err := crypto.PublicKeyToAddress(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}

	return &Account{
		Index:      index,
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Path:       path,
		Unlocked:   true, // 니모닉 기반이므로 언락됨
	}, nil
}

// 지갑 저장
func (p *WalletManager) SaveWallet() error {
	if p.wallet == nil {
		return fmt.Errorf("wallet not initialized")
	}

	// 지갑 디렉토리 생성
	if err := p.storage.CreateDir(p.walletDir); err != nil {
		return fmt.Errorf("failed to create wallet directory: %w", err)
	}

	// 파일명 생성 (wallet.json)
	walletFile := filepath.Join(p.walletDir, "wallet.json")

	bytes, err := json.Marshal(p.wallet)
	if err != nil {
		return fmt.Errorf("failed to marshal wallet data: %w", err)
	}

	err = p.storage.Write(walletFile, bytes)
	if err != nil {
		return fmt.Errorf("failed to write wallet file: %w", err)
	}

	return nil
}

// 지갑 로드
func (p *WalletManager) LoadWallet() error {
	// TODO: 지갑 파일에서 로드 구현
	return fmt.Errorf("not implemented yet")
}

// 현재 계정 가져오기
func (p *WalletManager) GetCurrentAccount() (*Account, error) {
	if p.wallet == nil || len(p.wallet.Accounts) == 0 {
		return nil, fmt.Errorf("no accounts available")
	}
	return p.wallet.Accounts[p.wallet.CurrentIndex], nil
}

// 계정 변경
func (p *WalletManager) SwitchAccount(index int) error {
	if p.wallet == nil {
		return fmt.Errorf("wallet not initialized")
	}

	if index < 0 || index >= len(p.wallet.Accounts) {
		return fmt.Errorf("invalid account index")
	}

	p.wallet.CurrentIndex = index
	return nil
}

// 니모닉 표시
func (p *WalletManager) GetMnemonic() (string, error) {
	if p.wallet == nil {
		return "", fmt.Errorf("wallet not initialized")
	}
	return p.wallet.Mnemonic, nil
}

// 계정 목록 가져오기
func (p *WalletManager) GetAccounts() ([]*Account, error) {
	if p.wallet == nil {
		return nil, fmt.Errorf("wallet not initialized")
	}
	return p.wallet.Accounts, nil
}
