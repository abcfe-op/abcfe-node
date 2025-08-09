package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"

	prt "github.com/abcfe/abcfe-node/protocol"
	"golang.org/x/crypto/sha3"
)

func PublicKeyToAddress(publicKey *ecdsa.PublicKey) (prt.Address, error) {
	// 공개키를 압축된 형태로 변환
	pubBytes := elliptic.MarshalCompressed(publicKey.Curve, publicKey.X, publicKey.Y)

	// Keccak256 해시
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubBytes[1:]) // 압축 접두사 제거
	hashBytes := hash.Sum(nil)

	// 마지막 20바이트를 Address 타입으로 변환
	var address prt.Address
	copy(address[:], hashBytes[len(hashBytes)-20:])

	return address, nil
}

// 주소에 0x 접두사 추가
func AddressTo0xPrefixString(address prt.Address) string {
	return "0x" + hex.EncodeToString(address[:])
}
