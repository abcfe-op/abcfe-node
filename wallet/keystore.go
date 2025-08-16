package wallet

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abcfe/abcfe-node/common/crypto"
	"github.com/abcfe/abcfe-node/config"
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
	Wallet    *MnemonicWallet
	Storage   WalletStorage
}

func NewWalletManager(walletDir string) *WalletManager {
	return &WalletManager{
		walletDir: walletDir,
		Storage:   FileWalletStorage{},
	}
}

func InitWallet(cfg *config.Config) (*WalletManager, error) {
	if cfg.Wallet.Path == "" {
		return nil, fmt.Errorf("failed to find wallet path. path is nil")
	}

	wm := NewWalletManager(cfg.Wallet.Path)
	err := wm.LoadWalletFile()
	if err != nil {
		return nil, fmt.Errorf("failed to load wallet: %w", err)
	}

	return wm, nil
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

	// 마스터 키를 바이트로 변환
	masterKeyBytes, err := crypto.PrivateKeyToBytes(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal master key: %w", err)
	}

	wallet := &MnemonicWallet{
		Mnemonic:     mnemonic,
		Seed:         seed,
		MasterKey:    masterKeyBytes,
		Accounts:     []*Account{},
		CurrentIndex: 0,
	}

	// 첫 번째 계정 생성
	account, err := p.deriveAccount(wallet, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive first account: %w", err)
	}

	wallet.Accounts = append(wallet.Accounts, account)
	p.Wallet = wallet

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

	// 마스터 키를 바이트로 변환
	masterKeyBytes, err := crypto.PrivateKeyToBytes(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal master key: %w", err)
	}

	wallet := &MnemonicWallet{
		Mnemonic:     mnemonic,
		Seed:         seed,
		MasterKey:    masterKeyBytes,
		Accounts:     []*Account{},
		CurrentIndex: 0,
	}

	// 첫 번째 계정 생성
	account, err := p.deriveAccount(wallet, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive first account: %w", err)
	}

	wallet.Accounts = append(wallet.Accounts, account)
	p.Wallet = wallet

	return wallet, nil
}

// 새 계정 추가
func (p *WalletManager) AddAccount() (*Account, error) {
	if p.Wallet == nil {
		return nil, fmt.Errorf("wallet not initialized")
	}

	nextIndex := len(p.Wallet.Accounts)
	account, err := p.deriveAccount(p.Wallet, nextIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account: %w", err)
	}

	p.Wallet.Accounts = append(p.Wallet.Accounts, account)
	return account, nil
}

// 계정 파생 (BIP-44 경로: m/44'/60'/0'/0/index)
func (p *WalletManager) deriveAccount(wallet *MnemonicWallet, index int) (*Account, error) {
	path := fmt.Sprintf("m/%d'/%d'/%d'/%d/%d",
		BIP44Purpose, BIP44CoinType, BIP44Account, BIP44Change, index)

	// 바이트를 개인키로 변환
	masterKey, err := crypto.BytesToPrivateKey(wallet.MasterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert master key: %w", err)
	}

	// 마스터 키에서 계정 키 파생
	privateKey, publicKey, err := crypto.DeriveAccountKey(masterKey, path)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account key: %w", err)
	}

	// 공개키에서 주소 생성
	address, err := crypto.PublicKeyToAddress(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}

	// 키들을 바이트로 변환
	privateKeyBytes, err := crypto.PrivateKeyToBytes(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert private key: %w", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert public key: %w", err)
	}

	return &Account{
		Index:      index,
		Address:    address,
		PrivateKey: privateKeyBytes,
		PublicKey:  publicKeyBytes,
		Path:       path,
		Unlocked:   true, // 니모닉 기반이므로 언락됨
	}, nil
}

// 지갑 저장
func (p *WalletManager) SaveWallet() error {
	if p.Wallet == nil {
		return fmt.Errorf("wallet not initialized")
	}

	// 지갑 디렉토리 생성
	if err := p.Storage.CreateDir(p.walletDir); err != nil {
		return fmt.Errorf("failed to create wallet directory: %w", err)
	}

	// JSON 직렬화
	bytes, err := json.Marshal(p.Wallet)
	if err != nil {
		return fmt.Errorf("failed to marshal wallet data: %w", err)
	}

	// 파일 저장
	walletFile := filepath.Join(p.walletDir, "wallet.json")
	err = p.Storage.Write(walletFile, bytes)
	if err != nil {
		return fmt.Errorf("failed to write wallet file: %w", err)
	}

	return nil
}

// 지갑 로드
func (p *WalletManager) LoadWalletFile() error {
	if p.walletDir == "" {
		return fmt.Errorf("wallet directory not set")
	}

	// 파일 경로 설정
	filePath := p.walletDir
	if !strings.HasSuffix(filePath, ".json") {
		filePath = filepath.Join(filePath, "wallet.json")
	}

	// 파일 읽기
	bytes, err := p.Storage.Read(filePath)
	if err != nil {
		return fmt.Errorf("failed to read wallet file: %w", err)
	}

	// JSON 역직렬화
	wallet := &MnemonicWallet{}
	err = json.Unmarshal(bytes, wallet)
	if err != nil {
		return fmt.Errorf("failed to unmarshal wallet data: %w", err)
	}

	p.Wallet = wallet
	return nil
}

// 현재 계정 가져오기
func (p *WalletManager) GetCurrentAccount() (*Account, error) {
	if p.Wallet == nil || len(p.Wallet.Accounts) == 0 {
		return nil, fmt.Errorf("no accounts available")
	}
	return p.Wallet.Accounts[p.Wallet.CurrentIndex], nil
}

// 계정 변경
func (p *WalletManager) SwitchAccount(index int) error {
	if p.Wallet == nil {
		return fmt.Errorf("Wallet not initialized")
	}

	if index < 0 || index >= len(p.Wallet.Accounts) {
		return fmt.Errorf("invalid account index")
	}

	p.Wallet.CurrentIndex = index
	return nil
}

// 니모닉 표시
func (p *WalletManager) GetMnemonic() (string, error) {
	if p.Wallet == nil {
		return "", fmt.Errorf("wallet not initialized")
	}
	return p.Wallet.Mnemonic, nil
}

// 계정 목록 가져오기
func (p *WalletManager) GetAccounts() ([]*Account, error) {
	if p.Wallet == nil {
		return nil, fmt.Errorf("wallet not initialized")
	}
	return p.Wallet.Accounts, nil
}

// 현재 계정의 개인키 가져오기 (ecdsa.PrivateKey 형태)
func (p *WalletManager) GetCurrentPrivateKey() (*ecdsa.PrivateKey, error) {
	if p.Wallet == nil || len(p.Wallet.Accounts) == 0 {
		return nil, fmt.Errorf("no accounts available")
	}

	account := p.Wallet.Accounts[p.Wallet.CurrentIndex]
	if len(account.PrivateKey) == 0 {
		return nil, fmt.Errorf("private key not available")
	}

	return crypto.BytesToPrivateKey(account.PrivateKey)
}

// 현재 계정의 공개키 가져오기 (ecdsa.PublicKey 형태)
func (p *WalletManager) GetCurrentPublicKey() (*ecdsa.PublicKey, error) {
	if p.Wallet == nil || len(p.Wallet.Accounts) == 0 {
		return nil, fmt.Errorf("no accounts available")
	}

	account := p.Wallet.Accounts[p.Wallet.CurrentIndex]
	if len(account.PublicKey) == 0 {
		return nil, fmt.Errorf("public key not available")
	}

	pub, err := x509.ParsePKIXPublicKey(account.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return pub.(*ecdsa.PublicKey), nil
}

// 특정 계정의 개인키 가져오기
func (p *WalletManager) GetAccountPrivateKey(index int) (*ecdsa.PrivateKey, error) {
	if p.Wallet == nil || index < 0 || index >= len(p.Wallet.Accounts) {
		return nil, fmt.Errorf("invalid account index")
	}

	account := p.Wallet.Accounts[index]
	if len(account.PrivateKey) == 0 {
		return nil, fmt.Errorf("private key not available")
	}

	return crypto.BytesToPrivateKey(account.PrivateKey)
}
