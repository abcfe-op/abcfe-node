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
