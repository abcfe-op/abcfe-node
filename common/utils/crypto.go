package utils

import (
	"crypto/sha256"
	"fmt"

	prt "github.com/abcfe/abcfe-node/protocol"
)

// 인터페이스를 가져와 해당 내용을 해싱한 후 해시의 16진수 인코딩을 반환
func Hash(i interface{}) prt.Hash {
	s := fmt.Sprintf("%v", i) // v: default formmater
	hash := sha256.Sum256([]byte(s))

	return hash
}
