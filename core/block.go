package core

import (
	"time"

	"github.com/abcfe/abcfe-node/common/utils"
	prt "github.com/abcfe/abcfe-node/protocol"
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
		Version:   p.cfg.Common.Version,
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

	block := &Block{
		Header:       *blkHeader,
		Transactions: txs,
	}

	blkHash := utils.Hash(block)
	block.Hash = blkHash

	return block
}

func GetBlock(height uint64) {

}
