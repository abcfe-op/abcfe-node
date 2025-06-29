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

// UTXO 관련 접두사
func (p *BlockChain) UpdateUtxo(batch *leveldb.Batch, blk Block) error {
	for _, tx := range blk.Transactions {
		if blk.Header.Height > 0 { // Genesis Block은 output만 처리
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

				// UTXO가 이미 사용되었는지 검증
				if utxo.Spent {
					return fmt.Errorf("UTXO already spent: %s:%d", utils.HashToString(input.TxID), input.OutputIndex)
				}

				utxo.Spent = true                    // utxo 사용처리
				utxo.SpentHeight = blk.Header.Height // 소비된 블록 높이

				utxoUpdBytes, err := utils.SerializeData(utxo, utils.SerializationFormatGob)
				if err != nil {
					return fmt.Errorf("failed to block serialization: %w", err)
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
		}

		// ! 아래 for 문에서 매번 키값을 가져오고 거기에 덮어쓰는 방향으로 코드가 진행 중.
		// 정확한 방식은 좀 소통하고 찾아봐야할 듯 함
		// // addr:[utxoKey1, utxoKey2]
		// addrList := []string{}

		// 3. OUTPUT으로 생성된 UTXO를 추가해줘야함
		for outputIndex, output := range tx.Outputs {
			newUtxo := UTXO{
				TxId:        tx.ID,
				OutputIndex: uint64(outputIndex),
				TxOut:       *output,
				Height:      blk.Header.Height,
				Spent:       false,
				SpentHeight: 0,
			}

			// UTXO 저장
			utxoKey := utils.GetUtxoKey(tx.ID, outputIndex)
			utxoBytes, err := utils.SerializeData(newUtxo, utils.SerializationFormatGob)
			if err != nil {
				return fmt.Errorf("failed to serialize utxo: %w", err)
			}
			batch.Put(utxoKey, utxoBytes)

			// 해당 주소의 UTXO 리스트에 추가
			var utxoList AddrUTXOSet
			utxoListKey := utils.GetUtxoListKey(output.Address)

			utxoListBytes, err := p.db.Get(utxoListKey, nil)
			if err == nil {
				if err := utils.DeserializeData(utxoListBytes, &utxoList, utils.SerializationFormatGob); err != nil {
					return fmt.Errorf("failed to deserialize utxo list: %w", err)
				}
			} else if err != leveldb.ErrNotFound {
				return fmt.Errorf("failed to get utxo list: %w", err)
			} else {
				utxoList = make(AddrUTXOSet)
			}

			// 리스트에 추가
			utxoList[string(utxoKey)] = true

			// 업데이트된 리스트 저장
			updatedListBytes, err := utils.SerializeData(utxoList, utils.SerializationFormatGob)
			if err != nil {
				return fmt.Errorf("failed to serialize utxo list: %w", err)
			}
			batch.Put(utxoListKey, updatedListBytes)
		}
	}

	return nil
}

// mempool 에 들어간 utox는 일단 사용한 것이라고 판단. 그 이후 블록 제안이 실패하면 그때ㄴ느 다시 사용할 수 있는 걸로.
// 즉 GetUtxoList에서 가져오는 Utxo에는 Mempool에 들어갔는지 체크하고 존재하면, 그 UTXO만 빼고 제공
func (p *BlockChain) GetUtxoList(address prt.Address) ([]*UTXO, error) {
	utxoListKey := utils.GetUtxoListKey(address)
	utxoListBytes, err := p.db.Get(utxoListKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get utxo data from db: %w", err)
	}

	var utxoList AddrUTXOSet
	if err := utils.DeserializeData(utxoListBytes, &utxoList, utils.SerializationFormatGob); err != nil {
		return nil, fmt.Errorf("failed to deserialize utxo list: %w", err)
	}

	var result []*UTXO
	for utxoKey := range utxoList {
		utxoBytes, err := p.db.Get([]byte(utxoKey), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get utxo data from db: %w", err)
		}

		var utxo UTXO
		if err := utils.DeserializeData(utxoBytes, &utxo, utils.SerializationFormatGob); err != nil {
			return nil, fmt.Errorf("failed to deserialize utxo list: %w", err)
		}

		// 멤풀에 들어간 UTXO인지 확인
		if p.isOnMempool(utxo.TxId, utxo.OutputIndex) {
			continue // 멤풀에 있는 UTXO // 사용됨으로 간주
		}

		result = append(result, &utxo)
	}

	return result, nil
}

func (p *BlockChain) CalBalanceUtxo(utxos []*UTXO) uint64 {
	var amount uint64
	for _, utxo := range utxos {
		if !utxo.Spent {
			amount += utxo.TxOut.Amount
		}
	}
	return amount
}
