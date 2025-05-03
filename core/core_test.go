package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/abcfe/abcfe-node/common/utils"
	prt "github.com/abcfe/abcfe-node/protocol"
)

// 트랜잭션 생성 헬퍼 함수
func setTestTransaction() *Transaction {
	// 입력 생성
	input := &TxInput{
		TxID:        prt.Hash{0x1, 0x2, 0x3}, // 임의의 해시
		OutputIndex: 0,
		Signature:   prt.Signature{}, // 빈 서명
		PublicKey:   []byte("Test PubKey"),
		// Sequence:    1,
	}

	// 출력 생성
	var addr prt.Address
	copy(addr[:], []byte("0x12300000000000000")) // 20바이트 주소

	output := &TxOutput{
		Address: addr,
		Amount:  1000,
		TxType:  0, // 일반 트랜잭션
	}

	// 트랜잭션 생성
	tx := &Transaction{
		Version:   "1.0.2",
		Timestamp: time.Now().Unix(),
		Inputs:    []*TxInput{input},
		Outputs:   []*TxOutput{output},
		Memo:      "Test Transaction",
		Data:      []byte{},
	}

	// ID 설정 (일반적으로는 트랜잭션 내용에 기반한 해시)
	tx.ID = utils.Hash(tx)

	return tx
}

// 멤풀 추가 삭제 조회 테스트
func TestAddTxToMempool(t *testing.T) {
	// mempool init
	mempool := NewMempool()

	// tx set
	tx := setTestTransaction()

	// tx add
	if err := mempool.NewTranaction(tx); err != nil {
		fmt.Println(err)
	}

	// tx get
	savedTx := mempool.GetTx(tx.ID)
	fmt.Println(savedTx)

	// tx delete
	mempool.DelTx(tx.ID)

	// tx get
	delCheckTx := mempool.GetTx(tx.ID)
	if delCheckTx != nil {
		fmt.Println("failed to delete saved tx")
	}
	fmt.Println("del done")
}

// 블록 구조체 생성 테스트
func TestSetBlock(t *testing.T) {
	// mempool init
	mempool := NewMempool()

	// tx add
	numTxs := 3
	for i := 0; i < numTxs; i++ {
		// tx set
		tx := setTestTransaction()

		// set unique tx.ID
		copy(tx.ID[:], []byte{byte(i), byte(i + 1), byte(i + 2)})

		// tx add to mempool
		if err := mempool.NewTranaction(tx); err != nil {
			fmt.Println(err)
		}
	}

	// test height
	height := uint64(1231)

	// test previous hash
	var prevHash prt.Hash
	copy(prevHash[:], []byte("0x1230000000000000012345678900")) // 32바이트 주소

	// set block header
	blkHeader := &BlockHeader{
		Version:   "1.0",
		Height:    height,
		PrevHash:  prevHash,
		Timestamp: time.Now().Unix(),
	}

	// set block
	// get tx from mempool
	txs := mempool.GetTxs()

	block := &Block{
		Header:       *blkHeader,
		Transactions: txs,
	}

	blkHash := utils.Hash(block)
	block.Hash = blkHash

	fmt.Println(block)
}

// func TestAddBlock(t *testing.T) {
// 	// mempool init
// 	mempool := NewMempool()

// 	// tx add
// 	numTxs := 3
// 	for i := 0; i < numTxs; i++ {
// 		// tx set
// 		tx := setTestTransaction()

// 		// set unique tx.ID
// 		copy(tx.ID[:], []byte{byte(i), byte(i + 1), byte(i + 2)})

// 		// tx add to mempool
// 		if err := mempool.NewTranaction(tx); err != nil {
// 			fmt.Println(err)
// 		}
// 	}

// 	// test height
// 	height := uint64(1231)

// 	// test previous hash
// 	var prevHash prt.Hash
// 	copy(prevHash[:], []byte("0x1230000000000000012345678900")) // 32바이트 주소

// 	// set block header
// 	blkHeader := &BlockHeader{
// 		Version:   "1.0",
// 		Height:    height,
// 		PrevHash:  prevHash,
// 		Timestamp: time.Now().Unix(),
// 	}

// 	// set block
// 	// get tx from mempool
// 	txs := mempool.GetTxs()

// 	block := &Block{
// 		Header:       *blkHeader,
// 		Transactions: txs,
// 	}

// 	blkHash := utils.Hash(block)
// 	block.Hash = blkHash

// 	// block serialization
// 	blockBytes, err := utils.SerializeData(block, utils.SerializationFormatGob)
// 	if err != nil {
// 		fmt.Println("here %w", err)
// 	}

// 	// db batch process ready
// 	batch := new(leveldb.Batch)

// 	// block hash - block data mapping
// 	blockHashKey := utils.GetBlockHashKey(prt.PrefixBlock, block.Hash)
// 	batch.Put(blockHashKey, blockBytes)

// 	// block height - block hash mapping
// 	heightKey := utils.GetBlockHeightKey(prt.PrefixBlockByHeight, block.Header.Height)
// 	batch.Put(heightKey, []byte(utils.HashToString(block.Hash)))

// 	db := chain.GetDB()

// 	// batch excute
// 	if err := p.db.Write(batch, nil); err != nil {
// 		return false, fmt.Errorf("failed to write batch: %w", err)
// 	}

// 	return true, nil
// }

func TestSetGenesisBlock(t *testing.T) {
	var defaultPrevHash prt.Hash

	for i := range defaultPrevHash {
		defaultPrevHash[i] = 0x00
	}

	blkHeader := &BlockHeader{
		Version:   "v0.1",
		Height:    0,
		PrevHash:  defaultPrevHash,
		Timestamp: time.Now().Unix(),
	}

	txIns := []*TxInput{}
	txOuts := []*TxOutput{}

	// 배열 초기화 - 올바른 문법으로 수정
	systemAddrs := []string{"ABCFEABCFEABCFEABCFEABCFEABCFEABCFEABCFE", "0000000000000000000000000000000000000000"}
	systemBals := []uint64{10000, 3300000}

	if len(systemAddrs) != len(systemBals) {
		fmt.Println("system address and balance count mismatch")
	}

	for i, systemAddr := range systemAddrs {
		addr, err := utils.StringToAddress(systemAddr)
		if err != nil {
			fmt.Println("failed to convert between address and string: ", err)
		}

		output := &TxOutput{
			Address: addr,
			Amount:  systemBals[i],
			TxType:  TxTypeGeneral,
		}
		txOuts = append(txOuts, output)
	}

	txs := []*Transaction{
		{
			Version:   "2.1",
			Timestamp: time.Now().Unix(),
			Inputs:    txIns,
			Outputs:   txOuts,
			Memo:      "ABCFE Chain Genesis Block",
		},
	}

	// TODO 서명 포함하고 그 이후 ID를 만들어야함

	for i, tx := range txs {
		txHash := utils.Hash(tx)
		txs[i].ID = txHash
	}

	block := &Block{
		Header:       *blkHeader,
		Transactions: txs,
	}
	block.Hash = utils.Hash(block)

	fmt.Println("Genesis Block: ", block)
}

// func setTestGenesisTxs() ([]*Transaction, error) {
// 	txIns := []*TxInput{}
// 	txOuts := []*TxOutput{}

// 	// 배열 초기화 - 올바른 문법으로 수정
// 	systemAddrs := []string{"0xABCFEABCFEABCFEABCFEABCFEABCFE"}
// 	systemBals := []uint64{10000}

// 	if len(systemAddrs) != len(systemBals) {
// 		return nil, fmt.Errorf("system address and balance count mismatch")

// 	}

// 	for i, systemAddr := range systemAddrs {
// 		addr, err := utils.StringToAddress(systemAddr)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to convert between address and string")
// 		}

// 		output := &TxOutput{
// 			Address: addr,
// 			Amount:  systemBals[i],
// 			TxType:  TxTypeGeneral,
// 		}
// 		txOuts = append(txOuts, output)
// 	}

// 	txs := []*Transaction{
// 		{
// 			Version:   "2.1",
// 			Timestamp: time.Now().Unix(),
// 			Inputs:    txIns,
// 			Outputs:   txOuts,
// 			Memo:      "ABCFE Chain Genesis Block",
// 		},
// 	}

// 	return txs, nil
// }
