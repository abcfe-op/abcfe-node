package protocol

import (
	"strconv"
)

const (
	// 블록체인 설정 정보
	PrefixNetworkConfig = "net:config"

	// 메타데이터 관련 접두사
	PrefixMeta          = "meta:"       // 메타데이터 키
	PrefixMetaHeight    = "meta:height" // 최신 블록 높이
	PrefixMetaBlockHash = "meta:hash"   // 최신 블록 해시

	// 블록 관련 접두사
	PrefixBlock         = "blk:"     // blk:해시 = 블록 데이터
	PrefixBlockByHeight = "blk:h:"   // blk:h:높이 = 블록 해시
	PrefixBlockTxs      = "blk:txs:" // blk:txs:블록해시:인덱스 = 트랜잭션 해시

	// 트랜잭션 관련 접두사
	// [사용패턴 1] tx:in:트랜잭션해시:인덱스 = 특정 input 데이터
	// [사용패턴 2] tx:in:트랜잭션해시 = 모든 input 데이터 목록
	PrefixTxIn = "tx:in:"

	// [사용패턴 1] tx:out:트랜잭션해시:인덱스 = 특정 output 데이터
	// [사용패턴 2] tx:out:트랜잭션해시 = 모든 output 데이터 목록
	PrefixTxOut    = "tx:out:"
	PrefixTxs      = "tx:"        // tx:트랜잭션해시 = 트랜잭션 데이터
	PrefixTxStatus = "tx:status:" // tx:status:해시 = 상태
	PrefixTxBlock  = "tx:blk:"    // tx:block:트랜잭션해시 = 블록해시

	// UTXO 관련 접두사
	PrefixUtxo          = "utxo:"      // utxo:트랜잭션해시:인덱스 = UTXO 데이터
	PrefixUtxoByAddress = "utxo:addr:" // utxo:addr:주소:트랜잭션해시:인덱스 = UTXO 데이터

	// 계정 관련 접두사
	PrefixAccount         = "acc:"      // acc:계정주소 = 계정 데이터
	PrefixAccountTxs      = "acc:txs:"  // acc:txs:계정주소 = 트랜잭션 해시 json-array
	PrefixAccountReceived = "acc:recv:" // acc:recv:계정주소:인덱스 = []{트랜잭션 해시: index} (수신)
	PrefixAccountSent     = "acc:sent:" // acc:sent:계정주소:인덱스 = []트랜잭션 해시 (발신)

	// 컨센서스 관련 접두사
	PrefixStakerInfo = "acc:staker:" // 지갑 주소 - 스테이킹 정보
)

// key gen
func GetTxInputKey(txHash string, index int) string {
	if index >= 0 { // index 유무로 특정 / 전체 구분
		return PrefixTxIn + txHash + ":" + strconv.Itoa(index)
	}
	return PrefixTxIn + txHash // 모든 입력 데이터 접근
}

func GetTxOutputKey(txHash string, index int) string {
	if index >= 0 { // index 유무로 특정 / 전체 구분
		return PrefixTxOut + txHash + ":" + strconv.Itoa(index)
	}
	return PrefixTxOut + txHash // 모든 출력 데이터 접근
}
