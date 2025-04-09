package utils

import (
	"encoding/hex"
	"fmt"

	prt "github.com/abcfe/abcfe-node/protocol"
)

// HashToString Hash 타입을 16진수 문자열로 변환
func HashToString(hash prt.Hash) string {
	return hex.EncodeToString(hash[:])
}

// StringToHash 16진수 문자열을 Hash 타입으로 변환
func StringToHash(str string) (prt.Hash, error) {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return prt.Hash{}, fmt.Errorf("잘못된 해시 문자열: %v", err)
	}

	// 해시 길이 검증
	if len(bytes) != 32 {
		return prt.Hash{}, fmt.Errorf("잘못된 해시 길이: %d (32바이트 필요)", len(bytes))
	}

	var hash prt.Hash
	copy(hash[:], bytes)
	return hash, nil
}

// BytesToHash 바이트 배열을 Hash 타입으로 변환
func BytesToHash(bytes []byte) prt.Hash {
	return prt.Hash(bytes)
}

// HashToBytes Hash 타입을 바이트 배열로 변환
func HashToBytes(hash prt.Hash) []byte {
	bytes := make([]byte, len(hash))
	copy(bytes, hash[:])
	return bytes
}

// AddressToString Address 타입을 16진수 문자열로 변환
func AddressToString(address prt.Address) string {
	return hex.EncodeToString(address[:])
}

// StringToAddress 16진수 문자열을 Address 타입으로 변환
func StringToAddress(str string) (prt.Address, error) {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return prt.Address{}, fmt.Errorf("잘못된 주소 문자열: %v", err)
	}

	// 주소 길이 검증
	if len(bytes) != 20 {
		return prt.Address{}, fmt.Errorf("잘못된 주소 길이: %d (20바이트 필요)", len(bytes))
	}

	var address prt.Address
	copy(address[:], bytes)
	return address, nil
}

// SignatureToString Signature 타입을 16진수 문자열로 변환
func SignatureToString(sig prt.Signature) string {
	return hex.EncodeToString(sig[:])
}

// StringToSignature 16진수 문자열을 Signature 타입으로 변환
func StringToSignature(str string) (prt.Signature, error) {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return prt.Signature{}, fmt.Errorf("잘못된 서명 문자열: %v", err)
	}

	// 서명 길이 검증
	if len(bytes) != 65 {
		return prt.Signature{}, fmt.Errorf("잘못된 서명 길이: %d (65바이트 필요)", len(bytes))
	}

	var sig prt.Signature
	copy(sig[:], bytes)
	return sig, nil
}
