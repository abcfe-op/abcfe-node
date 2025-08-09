package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/abcfe/abcfe-node/common/crypto"
	"github.com/abcfe/abcfe-node/common/utils"
	prt "github.com/abcfe/abcfe-node/protocol"
)

// 새 지갑 생성
func TestCreateWallet(t *testing.T) {
	w := Wallet{}

	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	w.PrivateKey = privateKey
	w.PublicKey = publicKey
	w.Address, err = crypto.PublicKeyToAddress(publicKey)
	if err != nil {
		t.Fatalf("failed to convert publicKey to address: %v", err)
	}

	// 올바른 출력 방법
	fmt.Printf("PrivateKey: %v\n", w.PrivateKey)
	fmt.Printf("PublicKey: %v\n", w.PublicKey)
	fmt.Printf("Address: %s\n", crypto.AddressTo0xPrefixString(w.Address)) // 0x 접두사 추가
}

// 트랜잭션 서명
func TestSignTransaction(t *testing.T) {
	w := Wallet{}

	privateKey, _, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	w.PrivateKey = privateKey

	// 테스트용 데이터
	testData := []byte("test transaction data")

	// ECDSA 서명
	signature, err := ecdsa.SignASN1(rand.Reader, w.PrivateKey, testData)
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}

	var sig prt.Signature
	copy(sig[:], signature)
	strSig := utils.SignatureToString(sig)
	fmt.Printf("Signature: %v\n", strSig)
}

// 서명 검증
func TestVerifySignature(t *testing.T) {
	w := Wallet{}

	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}

	w.PrivateKey = privateKey
	w.PublicKey = publicKey

	// 테스트용 데이터
	testData := []byte("test transaction data")

	// ECDSA 서명
	signature, err := ecdsa.SignASN1(rand.Reader, w.PrivateKey, testData)
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}

	var sig prt.Signature
	copy(sig[:], signature)
	strSig := utils.SignatureToString(sig)
	fmt.Printf("Signature: %v\n", strSig)

	// 원본 데이터로 검증 (수정됨!)
	result := ecdsa.VerifyASN1(w.PublicKey, testData, sig[:])
	fmt.Printf("Result: %t\n", result)
}

func TestKeyManager(t *testing.T) {
	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Error generating key pair: %v", err)
	}

	fmt.Printf("Private Key: %v\n", privateKey)
	fmt.Printf("Public Key: %v\n", publicKey)

	address, err := crypto.PublicKeyToAddress(publicKey)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	fmt.Printf("Address: %s\n", crypto.AddressTo0xPrefixString(address))
}
