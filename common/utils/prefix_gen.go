package utils

import (
	"strconv"

	prt "github.com/abcfe/abcfe-node/protocol"
)

// "blk:h:"
func GetBlockHeightKey(height uint64) []byte {
	hStr := Uint64ToString(height)
	hKey := []byte(prt.PrefixBlockByHeight + hStr)
	return hKey
}

// "blk:"
func GetBlockHashKey(hash prt.Hash) []byte {
	blkHashStr := HashToString(hash)
	blkKey := []byte(prt.PrefixBlock + blkHashStr)
	return blkKey
}

// "tx:"
func GetTxHashKey(txHash prt.Hash) []byte {
	txHashStr := HashToString(txHash)
	txKey := []byte(prt.PrefixTxs + txHashStr)
	return txKey
}

// "tx:status:"
func GetTxStatusKey(txHash prt.Hash) []byte {
	txHashStr := HashToString(txHash)
	txKey := []byte(prt.PrefixTxStatus + txHashStr)
	return txKey
}

// "tx:blk:"
func GetTxBlockHashKey(txHash prt.Hash) []byte {
	txHashStr := HashToString(txHash)
	txKey := []byte(prt.PrefixTxBlock + txHashStr)
	return txKey
}

// "tx:in:"
// [사용패턴 1] tx:in:트랜잭션해시:인덱스 = 특정 input 데이터
// [사용패턴 2] tx:in:트랜잭션해시 = 모든 input 데이터 목록
func GetTxInputKey(txHash prt.Hash, index int) []byte {
	txHashStr := HashToString(txHash)
	if index >= 0 { // index 유무로 특정 / 전체 구분
		return []byte(prt.PrefixTxIn + txHashStr + ":" + strconv.Itoa(index))
	}
	return []byte(prt.PrefixTxIn + txHashStr) // 모든 입력 데이터 접근
}

// "tx:out:"
// [사용패턴 1] tx:out:트랜잭션해시:인덱스 = 특정 output 데이터
// [사용패턴 2] tx:out:트랜잭션해시 = 모든 output 데이터 목록
func GetTxOutputKey(txHash prt.Hash, index int) []byte {
	txHashStr := HashToString(txHash)
	if index >= 0 { // index 유무로 특정 / 전체 구분
		return []byte(prt.PrefixTxOut + txHashStr + ":" + strconv.Itoa(index))
	}
	return []byte(prt.PrefixTxOut + txHashStr) // 모든 출력 데이터 접근
}
