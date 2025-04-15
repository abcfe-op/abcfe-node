package core

import (
	prt "github.com/abcfe/abcfe-node/protocol"
)

type UTXOSet map[string]*UTXO // Key: TxId + OutputIndex 문자열 조합

type UTXO struct {
	TxId        prt.Hash
	OutputIndex uint64
	TxOut       TxOutput
	Height      uint64
	Spent       bool // true : spent
}

// 제대로
// 그저 가져오는게 목적이라면 prefix 이용해서 DB 조회. 풀 노드 조회하면 너무 비효율적임
// 검증 목적이라면, txout 다 가져와서 utxo 구분해서 진행 필요.
func GetUtxos(address prt.Address) UTXOSet {
	var utxos UTXOSet

	// 해당 계정에 수신으로 생성된 트랜잭션 hash를 전부 가져오고
	recvTxHashes := getdb(prt.PrefixAccountReceived + address + ":")
	for _, txIdMap := range recvTxHashes { // txIdMap : []{hash: index}
		// prefix ins에 존재하는지 확인
		inputs := getdb(prt.PrefixTxIn + txIdMap.index + ":" + )
		for _, input := range inputs {
			// 없으면 out에 추가
		}
	}

}

func CalBalanceUtxo(utxos UTXOSet) uint64 {
	var amount uint64
	for _, utxo := range utxos {
		if !utxo.Spent {
			amount += utxo.TxOut.Amount
		}
	}
	return amount
}
