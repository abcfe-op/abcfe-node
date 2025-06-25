package core

import (
	"fmt"

	"github.com/abcfe/abcfe-node/common/utils"
	prt "github.com/abcfe/abcfe-node/protocol"
	"github.com/syndtr/goleveldb/leveldb"
)

type UTXOSet map[string]*UTXO    // Key: TxId + OutputIndex 문자열 조합
type AddrUTXOSet map[string]bool // user의 UTXO key 리스트
type UTXO struct {
	TxId        prt.Hash
	OutputIndex uint64
	TxOut       TxOutput
	Height      uint64
	Spent       bool // true : spent
	SpentHeight uint64
}

// // UTXO 관련 접두사
// PrefixUtxo        = "utxo:"      // utxo:트랜잭션해시:인덱스 = UTXO 데이터
// PrefixAddressUtxo = "utxo:addr:" // utxo:addr:주소 = UTXO 키 배열
// PrefixBalance     = "utxo:bal:"  // utxo:bal:주소 = 잔액
func (p *BlockChain) UpdateUtxo(batch *leveldb.Batch, blk Block) error {
	for _, tx := range blk.Transactions {
		for _, input := range tx.Inputs {
			// 1. input으로 사용한 UTXO 사용처리
			utxoKey := utils.GetUtxoKey(input.TxID, int(input.OutputIndex))
			utxoBytes, err := p.db.Get(utxoKey, nil)
			if err != nil {
				return fmt.Errorf("failed to get utxo data from db: %w", err)
			}

			var utxo UTXO
			if err := utils.DeserializeData(utxoBytes, &utxo, utils.SerializationFormatGob); err != nil {
				return fmt.Errorf("failed to deserialize utxo data: %w", err)
			}

			utxo.Spent = true                    // utxo 사용처리
			utxo.SpentHeight = blk.Header.Height // 소비된 블록 높이

			utxoUpdBytes, err := utils.SerializeData(utxo, utils.SerializationFormatGob)
			if err != nil {
				return fmt.Errorf("failed to block serialization.")
			}

			batch.Put(utxoKey, utxoUpdBytes)

			// 2. 주소별 UTXO 리스트에서 해당 UTXO 제거
			utxoListKey := utils.GetUtxoListKey(utxo.TxOut.Address)
			utxoListBytes, err := p.db.Get(utxoListKey, nil)
			if err != nil {
				return fmt.Errorf("failed to get utxo data from db: %w", err)
			}

			var utxoList AddrUTXOSet
			if err := utils.DeserializeData(utxoListBytes, &utxoList, utils.SerializationFormatGob); err != nil {
				return fmt.Errorf("failed to deserialize utxo list: %w", err)
			}

			// Map에서 제거. O(1)
			delete(utxoList, string(utxoKey))

			updatedListBytes, err := utils.SerializeData(utxoList, utils.SerializationFormatGob)
			if err != nil {
				return fmt.Errorf("failed to serialize utxo list: %w", err)
			}
			batch.Put(utxoListKey, updatedListBytes)
		}

		// 3. OUTPUT으로 생성된 UTXO를 추가해줘야함
		// - 각 트랜잭션 출력에 대해 반복
		// - 새로운 UTXO 객체 생성 (Spent=false, 현재 블록 높이)
		// - UTXO를 DB 키로 저장 (utxo:txid:outputindex)
		// - 주소별 UTXO 인덱스에 추가 (addr:address:utxo:txid:outputindex)
		// - 주소별 잔액 업데이트 (기존 잔액 + 새 출력 금액)
		// - 모든 변경사항을 batch에 추가
	}

	// 3. 제네시스 블록 특별 처리
	// - 제네시스 블록의 경우 입력이 없으므로 출력만 처리

	// 4. 오류 처리
	// - DB 조회 실패시 오류 반환
	// - 직렬화/역직렬화 실패시 오류 반환
	// - 존재하지 않는 UTXO 참조시 오류 반환
}

func (p *BlockChain) LoadUtxoData(batch *leveldb.Batch, blk Block) error {
	panic("not yet")
}

// 제대로
// 그저 가져오는게 목적이라면 prefix 이용해서 DB 조회. 풀 노드 조회하면 너무 비효율적임
// 검증 목적이라면, txout 다 가져와서 utxo 구분해서 진행 필요.
func GetUtxos(address prt.Address) UTXOSet {
	panic("not yet")
	// var utxos UTXOSet

	// // 해당 계정에 수신으로 생성된 트랜잭션 hash를 전부 가져오고
	// recvTxHashes := getdb(prt.PrefixAccountReceived + address + ":")
	// for _, txIdMap := range recvTxHashes { // txIdMap : []{hash: index}
	// 	// prefix ins에 존재하는지 확인
	// 	inputs := getdb(prt.PrefixTxIn + txIdMap.index + ":" + )
	// 	for _, input := range inputs {
	// 		// 없으면 out에 추가
	// 	}
	// }

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

func ValidateUtxo(tx *Transaction) (bool, error) {
	panic("not yet")
}
