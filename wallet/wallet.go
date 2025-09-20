package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"

	"github.com/abcfe/abcfe-node/common/crypto"
	"github.com/abcfe/abcfe-node/common/utils"
	"github.com/abcfe/abcfe-node/core"
	prt "github.com/abcfe/abcfe-node/protocol"
)

type Wallet struct {
	Address       prt.Address       // 20바이트 주소
	PrivateKey    *ecdsa.PrivateKey // ECDSA 개인키
	PublicKey     *ecdsa.PublicKey  // ECDSA 공개키
	WalletManager *WalletManager
}

func NewWallet(keystore *WalletManager) *Wallet {
	return &Wallet{
		WalletManager: keystore,
	}
}

// 새 지갑 생성
func (w *Wallet) CreateWallet() error {
	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	w.PrivateKey = privateKey
	w.PublicKey = publicKey
	w.Address, err = crypto.PublicKeyToAddress(publicKey)
	if err != nil {
		return fmt.Errorf("failed to convert publicKey to address: %w", err)
	}

	return nil
}

// 트랜잭션 서명
func (w *Wallet) SignTransaction(tx *core.Transaction) (*prt.Signature, error) {
	if w.PrivateKey == nil {
		return nil, fmt.Errorf("wallet not unlocked")
	}

	// 트랜잭션 해시 계산
	txHash := utils.Hash(tx)
	txHashBytes := utils.HashToBytes(txHash)

	// ECDSA 서명
	signature, err := ecdsa.SignASN1(rand.Reader, w.PrivateKey, txHashBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	var sig prt.Signature
	copy(sig[:], signature)

	return &sig, nil
}

// 서명 검증
func (w *Wallet) VerifySignature(tx *core.Transaction, sig *prt.Signature, address prt.Address) bool {
	// 트랜잭션 해시 계산
	txHash := utils.Hash(tx)
	txHashBytes := utils.HashToBytes(txHash)

	// ECDSA 서명 검증
	return ecdsa.VerifyASN1(w.PublicKey, txHashBytes, sig[:])
}

// 트랜잭션 생성
func (w *Wallet) CreateTransaction() (*core.Transaction, error) {
	// TODO: 트랜잭션 생성 로직 구현
	return nil, fmt.Errorf("not implemented yet")
}
