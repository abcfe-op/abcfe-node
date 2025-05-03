package core

import (
	"fmt"
	"time"

	"github.com/abcfe/abcfe-node/common/utils"
	prt "github.com/abcfe/abcfe-node/protocol"
	"github.com/syndtr/goleveldb/leveldb"
)

type Block struct {
	Header       BlockHeader    `json:"header"`       // 블록 헤더
	Hash         prt.Hash       `json:"hash"`         // 블록 해시
	Transactions []*Transaction `json:"transactions"` // 트랜잭션 목록

	// TODO: Consensus 시작시
	// Proposer      Address     `json:"proposer"`      // 블록 제안자
	// Validators    []Address   `json:"validators"`    // 검증자 목록 // 자세한 정보는 상위에 있는 컨센서스 패키지에서 처리할 예정
	// Signatures    []Signature `json:"signatures"`    // 검증자 서명
	// ConsensusData []byte      `json:"consensusData"` // 컨센서스 관련 데이터는 단방향 참조를 위해 직렬화만
}

type BlockHeader struct {
	Version  string   `json:"version"`  // 블록체인 프로토콜 버전
	Height   uint64   `json:"height"`   // 블록 높이 (uint64로 변경)
	PrevHash prt.Hash `json:"prevHash"` // 이전 블록 해시
	// TODO: 고도화 시작시
	// MerkleRoot Hash   `json:"merkleRoot"` // 트랜잭션 머클 루트
	// StateRoot  Hash   `json:"stateRoot"`  // 상태 머클 루트 (UTXO 또는 계정 상태)
	Timestamp int64 `json:"timestamp"` // 블록 생성 시간 (Unix 타임스탬프)
}

func (p *BlockChain) setBlockHeader(height uint64, prevHash prt.Hash) *BlockHeader {
	blkHeader := &BlockHeader{
		Version:   p.cfg.Version.Protocol,
		Height:    height,
		PrevHash:  prevHash,
		Timestamp: time.Now().Unix(),
	}
	return blkHeader
}

func (p *BlockChain) SetBlock(prevHash prt.Hash, height uint64) *Block {
	blkHeader := p.setBlockHeader(height, prevHash)

	// 메모리 풀에서 트랜잭션 가져오기
	txs := p.mempool.GetTxs()

	blk := &Block{
		Header:       *blkHeader,
		Transactions: txs,
	}

	blkHash := utils.Hash(blk)
	blk.Hash = blkHash

	return blk
}

func (p *BlockChain) AddBlock(blk Block) (bool, error) {
	// block serialization
	blkBytes, err := utils.SerializeData(blk, utils.SerializationFormatGob)
	if err != nil {
		return false, fmt.Errorf("failed to block serialization.")
	}

	// db batch process ready
	batch := new(leveldb.Batch)

	// block hash - block data mapping
	blkHashKey := utils.GetBlockHashKey(prt.PrefixBlock, blk.Hash)
	batch.Put(blkHashKey, blkBytes)

	// block height - block hash mapping
	heightKey := utils.GetBlockHeightKey(prt.PrefixBlockByHeight, blk.Header.Height)
	batch.Put(heightKey, []byte(utils.HashToString(blk.Hash)))

	// TODO: transaction mapping

	// chain status update
	if blk.Header.Height > p.LatestHeight {
		if err := p.UpdateChainState(blk.Header.Height, utils.HashToString(blk.Hash)); err != nil {
			return false, fmt.Errorf("failed to update chain status: %w", err)
		}
	}

	// batch excute
	if err := p.db.Write(batch, nil); err != nil {
		return false, fmt.Errorf("failed to write batch: %w", err)
	}

	return true, nil
}

func (p *BlockChain) GetBlock(height uint64) (*Block, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	heightKey := utils.GetBlockHeightKey(prt.PrefixBlockByHeight, height)
	blkBytes, err := p.db.Get(heightKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get block data from db: %w", err)
	}

	// block data deserialization
	var block Block
	if err := utils.DeserializeData(blkBytes, &block, utils.SerializationFormatGob); err != nil {
		return nil, fmt.Errorf("failed to deserialize block data: %w", err)
	}

	return &block, nil
}

func (p *BlockChain) ValidateBlock(block Block) (bool, error) {
	panic("Not Developed Yet")
}

// BlockToJSON 블록을 JSON 형식으로 변환
func blockToJSON(block interface{}) ([]byte, error) {
	return utils.SerializeData(block, utils.SerializationFormatJSON)
}

// JSONToBlock JSON 형식에서 블록으로 변환
func jsonToBlock(data []byte, block interface{}) error {
	return utils.DeserializeData(data, block, utils.SerializationFormatJSON)
}
