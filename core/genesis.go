package core

import (
	"fmt"
	"time"

	"github.com/abcfe/abcfe-node/common/utils"
	prt "github.com/abcfe/abcfe-node/protocol"
)

func (p *BlockChain) SetGenesisBlock() (*Block, error) {
	var defaultPrevHash prt.Hash

	// 제네시스블록에는 이전 블록 해시값이 없으므로, 0으로 구성
	for i := range defaultPrevHash {
		defaultPrevHash[i] = 0x00
	}

	txs, err := p.setGenesisTxs()
	if err != nil {
		return nil, err
	}

	// 머클 루트 계산
	merkleRoot := calculateMerkleRoot(txs)

	blkHeader := &BlockHeader{
		PrevHash:   defaultPrevHash,
		Version:    p.cfg.Version.Protocol,
		Height:     0,
		Timestamp:  time.Now().Unix(),
		MerkleRoot: merkleRoot,
	}

	block := &Block{
		Header:       *blkHeader,
		Transactions: txs,
	}

	blkHash := utils.Hash(block)
	block.Header.Hash = blkHash

	return block, nil
}

func (p *BlockChain) setGenesisTxs() ([]*Transaction, error) {
	txIns := []*TxInput{}
	txOuts := []*TxOutput{}

	systemAddrs := p.cfg.Genesis.SystemAddresses
	systemBals := p.cfg.Genesis.SystemBalances

	if len(systemAddrs) != len(systemBals) {
		return nil, fmt.Errorf("system address and balance count mismatch")
	}

	for i, systemAddr := range systemAddrs {
		addr, err := utils.StringToAddress(systemAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert between address and string")
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
			Version:   p.cfg.Version.Transaction,
			Timestamp: time.Now().Unix(),
			Inputs:    txIns,
			Outputs:   txOuts,
			Memo:      "ABCFE Chain Genesis Block",
		},
	}

	// TODO 서명값 넣고 이후 해시화

	for i, tx := range txs {
		txHash := utils.Hash(tx)
		txs[i].ID = txHash
	}

	return txs, nil
}
