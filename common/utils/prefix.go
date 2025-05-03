package utils

import (
	prt "github.com/abcfe/abcfe-node/protocol"
)

func GetBlockHeightKey(prefix string, height uint64) []byte {
	hStr := Uint64ToString(height)
	hKey := []byte(prefix + hStr)
	return hKey
}

func GetBlockHashKey(prefix string, hash prt.Hash) []byte {
	blkHashStr := HashToString(hash)
	blkKey := []byte(prefix + blkHashStr)
	return blkKey
}
