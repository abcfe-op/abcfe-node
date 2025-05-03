package core

import (
	"fmt"
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
	Outputs []*TxOutput `json:"outputs"` // 트랜잭션 출력 // 잔돈 때문에 단일 출력이 아닌 다중 출력이어야함

	Memo string `json:"memo"` // 트랜잭션 메모 (inputData 대체)
	Data []byte `json:"data"` // 임의 데이터 (스마트 컨트랙트 호출 등)
	// Signature // TODO
}

type TxInput struct {
	TxID        prt.Hash      `json:"txId"`        // 참조 트랜잭션 ID
	OutputIndex uint64        `json:"outputIndex"` // 참조 출력 인덱스
	Signature   prt.Signature `json:"signature"`   // 서명
	PublicKey   []byte        `json:"publicKey"`   // 공개키
	// Sequence    uint64        `json:"sequence"`    // 시퀀스 번호 (RBF 지원)
}

type TxOutput struct {
	Address prt.Address `json:"address"` // 수신자 주소
	Amount  uint64      `json:"amount"`  // 금액 (int에서 uint64로 변경)
	TxType  uint8       `json:"txType"`  // 스크립트 타입 (일반/스테이킹/기타)
}

// Tx Input and Output pair
type TxIOPair struct {
	TxIns  []*TxInput  `json:"txIns"`
	TxOuts []*TxOutput `json:"txOuts"`
}

func (p *BlockChain) SetTx(from prt.Address, to prt.Address, amount uint64, memo string, data []byte, txType uint8) (*Transaction, error) {
	utxos := GetUtxos(from)
	if CalBalanceUtxo(utxos) < amount {
		return &Transaction{}, fmt.Errorf("not enough balance")
	}

	txInAndOut, err := p.setTxIOPair(utxos, from, to, amount, txType)
	if err != nil {
		return &Transaction{}, err
	}

	tx := Transaction{
		Version:   p.cfg.Version.Transaction,
		Timestamp: time.Now().Unix(),
		Inputs:    txInAndOut.TxIns,
		Outputs:   txInAndOut.TxOuts,
		Memo:      memo,
		Data:      data,
	}

	txHash := utils.Hash(tx)
	tx.ID = txHash

	return &tx, nil
}

// tx input과 output을 구성
func (p *BlockChain) setTxIOPair(utxos UTXOSet, from prt.Address, to prt.Address, amount uint64, txType uint8) (TxIOPair, error) {
	var txInAndOut TxIOPair
	var total uint64

	// set tx in
	for _, utxo := range utxos {
		// ! 이 부분 등호 이해 안됨
		if total >= amount {
			break
		}
		// ! pubkey는 전처리 시그니처 후처리 필요
		txIn := p.setTxInput(utxo.TxId, utxo.OutputIndex, prt.Signature{}, nil)
		txInAndOut.TxIns = append(txInAndOut.TxIns, txIn)

		total += utxo.TxOut.Amount
	}

	// utxo 탈탈 털었는데도 돈이 부족하다면 에러
	if total < amount {
		return TxIOPair{}, fmt.Errorf("Not enough balance: required %d, balance %d", amount, total)
	}

	// set tx out
	txOut := p.setTxOutput(to, amount, txType)
	txInAndOut.TxOuts = append(txInAndOut.TxOuts, txOut)

	// 거슬러줘야한다면 잔액 반환
	if total > amount {
		changeOut := total - amount
		txOut := p.setTxOutput(from, changeOut, txType)
		txInAndOut.TxOuts = append(txInAndOut.TxOuts, txOut)
	}

	return txInAndOut, nil
}

func (p *BlockChain) setTxInput(txOutID prt.Hash, txOutIdx uint64, sig prt.Signature, pubKey []byte) *TxInput {
	txIn := &TxInput{
		TxID:        txOutID,
		OutputIndex: txOutIdx,
		Signature:   sig,
		PublicKey:   pubKey,
	}

	return txIn
}

func (p *BlockChain) setTxOutput(toAddr prt.Address, amount uint64, txType uint8) *TxOutput {
	txOut := &TxOutput{
		Address: toAddr,
		Amount:  amount,
		TxType:  txType,
	}

	return txOut
}

func (p *BlockChain) GetTx(height uint64, id prt.Hash) {
	// prt.PrefixTx
}

func (p *BlockChain) GetTxStatus(height uint64, id prt.Hash) {
	// prt.PrefixTxStatus
}

func (p *BlockChain) GetTxs(height uint64) {
	// prt.PrefixTx
}
