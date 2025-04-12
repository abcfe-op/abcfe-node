package core

import (
	"time"

	"github.com/abcfe/abcfe-node/common/utils"
	prt "github.com/abcfe/abcfe-node/protocol"
)

// 트랜잭션 구조체
type Transaction struct {
	Version   string   `json:"version"`   // 트랜잭션 버전
	ID        prt.Hash `json:"id"`        // 트랜잭션 ID (해시)
	Timestamp int64    `json:"timestamp"` // 트랜잭션 생성 시간

	Inputs  []*TxInput  `json:"inputs"`  // 트랜잭션 입력
	Outputs []*TxOutput `json:"outputs"` // 트랜잭션 출력

	Memo string `json:"memo"` // 트랜잭션 메모 (inputData 대체)
	Data []byte `json:"data"` // 임의 데이터 (스마트 컨트랙트 호출 등)
}

type TxInput struct {
	TxID        prt.Hash      `json:"txId"`        // 참조 트랜잭션 ID
	OutputIndex uint64        `json:"outputIndex"` // 참조 출력 인덱스
	Signature   prt.Signature `json:"signature"`   // 서명
	PublicKey   []byte        `json:"publicKey"`   // 공개키
	Sequence    uint64        `json:"sequence"`    // 시퀀스 번호 (RBF 지원)
}

type TxOutput struct {
	Address prt.Address `json:"address"` // 수신자 주소
	Amount  uint64      `json:"amount"`  // 금액 (int에서 uint64로 변경)
	TxType  uint8       `json:"txType"`  // 스크립트 타입 (일반/스테이킹/기타)
}

func (p *BlockChain) SetTx(txIn []*TxInput, txOut []*TxOutput, memo string, data []byte) {
	tx := Transaction{
		Version:   p.cfg.Version.Transaction,
		Timestamp: time.Now().Unix(),
		Inputs:    txIn,
		Outputs:   txOut,
		Memo:      memo,
		Data:      data,
	}

	txHash := utils.Hash(tx)
	tx.ID = txHash

}

func (p *BlockChain) SetTxInput(txID prt.Hash, outIndex uint64, sig prt.Signature, pubKey []byte, seq uint64) *TxInput {
	txIn := &TxInput{
		TxID:        txID,
		OutputIndex: outIndex,
		Signature:   sig,
		PublicKey:   pubKey,
		Sequence:    seq,
	}

	return txIn
}

func (p *BlockChain) SetTxOutput(toAddr prt.Address, amount uint64, txType uint8) *TxOutput {
	txOut := &TxOutput{
		Address: toAddr,
		Amount:  amount,
		TxType:  txType,
	}

	return txOut
}

func (p *BlockChain) GetTx(height uint64, id prt.Hash) {}

func (p *BlockChain) GetTxStatus(height uint64, id prt.Hash) {}

func (p *BlockChain) GetTxs(height uint64) {}
