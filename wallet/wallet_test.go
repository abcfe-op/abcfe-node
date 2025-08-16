package wallet

import (
	"fmt"
	"testing"

	"github.com/abcfe/abcfe-node/common/crypto"
	"github.com/abcfe/abcfe-node/common/utils"
)

// 기존 테스트

// // 새 지갑 생성
// func TestCreateWallet(t *testing.T) {
// 	w := Wallet{}

// 	privateKey, publicKey, err := crypto.GenerateKeyPair()
// 	if err != nil {
// 		t.Fatalf("failed to generate key pair: %v", err)
// 	}

// 	w.PrivateKey = privateKey
// 	w.PublicKey = publicKey
// 	w.Address, err = crypto.PublicKeyToAddress(publicKey)
// 	if err != nil {
// 		t.Fatalf("failed to convert publicKey to address: %v", err)
// 	}

// 	// 올바른 출력 방법
// 	fmt.Printf("PrivateKey: %v\n", w.PrivateKey)
// 	fmt.Printf("PublicKey: %v\n", w.PublicKey)
// 	fmt.Printf("Address: %s\n", crypto.AddressTo0xPrefixString(w.Address)) // 0x 접두사 추가
// }

// // 트랜잭션 서명
// func TestSignTransaction(t *testing.T) {
// 	w := Wallet{}

// 	privateKey, _, err := crypto.GenerateKeyPair()
// 	if err != nil {
// 		t.Fatalf("failed to generate key pair: %v", err)
// 	}

// 	w.PrivateKey = privateKey

// 	// 테스트용 데이터
// 	testData := []byte("test transaction data")

// 	// ECDSA 서명
// 	signature, err := ecdsa.SignASN1(rand.Reader, w.PrivateKey, testData)
// 	if err != nil {
// 		t.Fatalf("failed to sign transaction: %v", err)
// 	}

// 	var sig prt.Signature
// 	copy(sig[:], signature)
// 	strSig := utils.SignatureToString(sig)
// 	fmt.Printf("Signature: %v\n", strSig)
// }

// // 서명 검증
// func TestVerifySignature(t *testing.T) {
// 	w := Wallet{}

// 	privateKey, publicKey, err := crypto.GenerateKeyPair()
// 	if err != nil {
// 		t.Fatalf("failed to generate key pair: %v", err)
// 	}

// 	w.PrivateKey = privateKey
// 	w.PublicKey = publicKey

// 	// 테스트용 데이터
// 	testData := []byte("test transaction data")

// 	// ECDSA 서명
// 	signature, err := ecdsa.SignASN1(rand.Reader, w.PrivateKey, testData)
// 	if err != nil {
// 		t.Fatalf("failed to sign transaction: %v", err)
// 	}

// 	var sig prt.Signature
// 	copy(sig[:], signature)
// 	strSig := utils.SignatureToString(sig)
// 	fmt.Printf("Signature: %v\n", strSig)

// 	// 원본 데이터로 검증 (수정됨!)
// 	result := ecdsa.VerifyASN1(w.PublicKey, testData, sig[:])
// 	fmt.Printf("Result: %t\n", result)
// }

// func TestKeyManager(t *testing.T) {
// 	privateKey, publicKey, err := crypto.GenerateKeyPair()
// 	if err != nil {
// 		t.Fatalf("Error generating key pair: %v", err)
// 	}

// 	fmt.Printf("Private Key: %v\n", privateKey)
// 	fmt.Printf("Public Key: %v\n", publicKey)

// 	address, err := crypto.PublicKeyToAddress(publicKey)
// 	if err != nil {
// 		t.Fatalf("Error: %v", err)
// 	}
// 	fmt.Printf("Address: %s\n", crypto.AddressTo0xPrefixString(address))
// }

// 지갑 생성
func TestCreateWallet(t *testing.T) {
	wm := NewWalletManager("./resource/wallet")
	wallet, err := wm.CreateWallet()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	wm.Wallet = wallet

	fmt.Printf("dir: %v\n", wm.walletDir)
	fmt.Printf("Wallet: %v\n", wm.Wallet.Mnemonic)
	fmt.Printf("seed: %v\n", wm.Wallet.Seed)
	fmt.Printf("master key: %v\n", wm.Wallet.MasterKey)
	fmt.Printf("accounts: %v\n", wm.Wallet.Accounts[0])
	fmt.Printf("pub key: %v\n", wm.Wallet.Accounts[0].PublicKey)
	fmt.Printf("priv key: %v\n", wm.Wallet.Accounts[0].PrivateKey)
	fmt.Printf("address: %v\n", crypto.AddressTo0xPrefixString(wm.Wallet.Accounts[0].Address))
	fmt.Printf("path: %v\n", wm.Wallet.Accounts[0].Path)
	fmt.Printf("cur idx: %v\n", wm.Wallet.CurrentIndex)
}

func TestRestoreWallet(t *testing.T) {
	wm := NewWalletManager("./resource/wallet")
	wallet, err := wm.CreateWallet()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	wm.Wallet = wallet

	newAddress := crypto.AddressTo0xPrefixString(wm.Wallet.Accounts[0].Address)

	// 복구 시작
	restoredWallet, err := wm.RestoreWallet(wm.Wallet.Mnemonic)
	if err != nil {
		t.Fatalf("Error: %v\n", err)
	}
	wm.Wallet = restoredWallet

	restoredAddress := crypto.AddressTo0xPrefixString(restoredWallet.Accounts[0].Address)
	if restoredAddress != newAddress {
		t.Fatalf("address is different each other. %s | %s", newAddress, restoredAddress)
	}

	fmt.Printf("first address: %s\n", newAddress)
	fmt.Printf("restored address: %s\n", restoredAddress)
}

func TestSaveWallet(t *testing.T) {
	// 디렉토리 경로로 변경
	wm := NewWalletManager("../resource/wallet")
	wallet, err := wm.CreateWallet()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	wm.Wallet = wallet

	err = wm.SaveWallet()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	fmt.Printf("지갑이 성공적으로 저장되었습니다: %s\n", wm.walletDir)
}

func TestLoadWallet(t *testing.T) {
	// 디렉토리 경로로 변경
	wm := NewWalletManager("../resource/wallet") // 파일 경로가 아닌 디렉토리 경로
	err := wm.LoadWalletFile()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	fmt.Printf("지갑이 성공적으로 로드되었습니다: %v\n", utils.AddressToString(wm.Wallet.Accounts[0].Address))
}
