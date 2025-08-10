package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

func GenerateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := privateKey.PublicKey
	return privateKey, &publicKey, err
}

// 시드에서 마스터 키 파생 (간단한 버전)
func DeriveMasterKey(seed []byte) (*ecdsa.PrivateKey, error) {
	// 시드를 SHA256으로 해시하여 개인키로 사용
	hash := sha256.Sum256(seed)

	// 해시를 개인키로 변환
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(hash[:])

	// 공개키 계산
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(hash[:])

	return privateKey, nil
}

// 경로에서 계정 키 파생 (간단한 버전)
func DeriveAccountKey(masterKey *ecdsa.PrivateKey, path string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	// 경로를 해시하여 계정별 고유한 오프셋 생성
	pathHash := sha256.Sum256([]byte(path))

	// 마스터 키에 오프셋 추가
	offset := new(big.Int).SetBytes(pathHash[:])
	newD := new(big.Int).Add(masterKey.D, offset)
	newD.Mod(newD, masterKey.PublicKey.Curve.Params().N)

	// 새 개인키 생성
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.D = newD

	// 공개키 계산
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(newD.Bytes())

	return privateKey, &privateKey.PublicKey, nil
}
