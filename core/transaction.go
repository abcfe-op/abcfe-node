package core

import (
	"crypto/sha256"
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
	// Status // TODO
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

func (p *BlockChain) SetTransferTx(from prt.Address, to prt.Address, amount uint64, memo string, data []byte, txType uint8) (*Transaction, error) {
	utxos, err := p.GetUtxoList(from)
	if err != nil {
		return nil, err
	}
	if p.CalBalanceUtxo(utxos) < amount {
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
func (p *BlockChain) setTxIOPair(utxos []*UTXO, from prt.Address, to prt.Address, amount uint64, txType uint8) (*TxIOPair, error) {
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
		return nil, fmt.Errorf("Not enough balance: required %d, balance %d", amount, total)
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

	return &txInAndOut, nil
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

func (p *BlockChain) GetTx(txId prt.Hash) (*Transaction, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := utils.GetTxHashKey(txId)
	txBytes, err := p.db.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get block hash from db: %w", err)
	}

	// tx data bytes -> tx data deserialization
	var tx Transaction
	if err := utils.DeserializeData(txBytes, &tx, utils.SerializationFormatGob); err != nil {
		return nil, fmt.Errorf("failed to deserialize block data: %w", err)
	}

	return &tx, nil
}

func (p *BlockChain) GetBlockHashByTxId(txId prt.Hash) (prt.Hash, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// block height -> block hash bytes
	key := utils.GetTxBlockHashKey(txId)
	blkHashBytes, err := p.db.Get(key, nil)
	if err != nil {
		return prt.Hash{}, fmt.Errorf("failed to get block hash from db: %w", err)
	}

	// block hash bytes -> block hash string
	var blkHash prt.Hash
	blkHash = utils.BytesToHash(blkHashBytes)

	return blkHash, nil
}

// TODO
// func (p *BlockChain) GetTxStatus(height uint64, id prt.Hash) {
// prt.PrefixTxStatus
// }

func (p *BlockChain) GetInputTx(txId prt.Hash) (*[]TxInput, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := utils.GetTxInputKey(txId, prt.WholeTxIdx)
	txBytes, err := p.db.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx input data from db: %w", err)
	}

	// tx data bytes -> tx data deserialization
	var txInputs []TxInput
	if err := utils.DeserializeData(txBytes, &txInputs, utils.SerializationFormatGob); err != nil {
		return nil, fmt.Errorf("failed to deserialize tx input data: %w", err)
	}

	return &txInputs, nil
}

func (p *BlockChain) GetOutputTx(txId prt.Hash) (*[]TxOutput, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := utils.GetTxOutputKey(txId, prt.WholeTxIdx)
	txBytes, err := p.db.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx output data from db: %w", err)
	}

	// tx data bytes -> tx data deserialization
	var txOutputs []TxOutput
	if err := utils.DeserializeData(txBytes, &txOutputs, utils.SerializationFormatGob); err != nil {
		return nil, fmt.Errorf("failed to deserialize tx output data: %w", err)
	}

	return &txOutputs, nil
}

func (p *BlockChain) GetInputTxByIdx(txId prt.Hash, idx int) (*TxInput, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := utils.GetTxInputKey(txId, idx)
	txBytes, err := p.db.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx input data from db: %w", err)
	}

	// tx data bytes -> tx data deserialization
	var txInput TxInput
	if err := utils.DeserializeData(txBytes, &txInput, utils.SerializationFormatGob); err != nil {
		return nil, fmt.Errorf("failed to deserialize tx input data: %w", err)
	}

	return &txInput, nil
}

func (p *BlockChain) GetOutputTxByIdx(txId prt.Hash, idx int) (*TxOutput, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := utils.GetTxOutputKey(txId, idx)
	txBytes, err := p.db.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx output data from db: %w", err)
	}

	// tx data bytes -> tx data deserialization
	var txOutput TxOutput
	if err := utils.DeserializeData(txBytes, &txOutput, utils.SerializationFormatGob); err != nil {
		return nil, fmt.Errorf("failed to deserialize tx output data: %w", err)
	}

	return &txOutput, nil
}

func calculateMerkleRoot(txs []*Transaction) prt.Hash {
	if len(txs) == 0 {
		return prt.Hash{} // 빈 해시 반환
	}

	// 각 트랜잭션을 해시
	hashes := make([]prt.Hash, len(txs))
	for i, tx := range txs {
		hashes[i] = utils.Hash(tx)
	}

	// 머클 트리 계산
	return buildMerkleTree(hashes)
}

func buildMerkleTree(hashes []prt.Hash) prt.Hash {
	if len(hashes) == 1 {
		return hashes[0]
	}

	// 짝수 개로 맞추기
	if len(hashes)%2 != 0 {
		hashes = append(hashes, hashes[len(hashes)-1])
	}

	// 다음 레벨 계산
	nextLevel := make([]prt.Hash, len(hashes)/2)
	for i := 0; i < len(hashes); i += 2 {
		combined := append(hashes[i][:], hashes[i+1][:]...)
		nextLevel[i/2] = sha256.Sum256(combined)
	}

	return buildMerkleTree(nextLevel)
}
