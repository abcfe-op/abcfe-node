package core

import (
	"fmt"
	"sync"

	prt "github.com/abcfe/abcfe-node/protocol"

	"github.com/abcfe/abcfe-node/common/utils"
)

type Mempool struct {
	transactions map[string]*Transaction
	mu           sync.RWMutex
}

func NewMempool() *Mempool {
	return &Mempool{
		transactions: make(map[string]*Transaction),
	}
}

func (p *Mempool) NewTranaction(tx *Transaction) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	txId := utils.HashToString(tx.ID)

	if _, exists := p.transactions[txId]; exists {
		return fmt.Errorf("tx already exists in mempool.")
	}

	// 트랜잭션이 이미 검증되었고 여기에 들어왔다고 판단? 아니면 추가 검증?
	// ValidateTx(tx)

	p.transactions[txId] = tx
	return nil
}

func (p *Mempool) GetTx(txId prt.Hash) *Transaction {
	p.mu.Lock()
	defer p.mu.Unlock()

	strTxId := utils.HashToString(txId)

	return p.transactions[strTxId]
}

// 블록 추가시 멤풀에서 트랜잭션 추출
func (p *Mempool) GetTxs() []*Transaction {
	p.mu.Lock()
	defer p.mu.Unlock()

	txs := make([]*Transaction, 0, len(p.transactions))
	for _, tx := range p.transactions {
		txs = append(txs, tx)
	}

	// TODO 가스비 로직이 추가되면, 가스비 기준으로 정렬

	// 트랜잭션 수가 최대값보다 많으면 제한
	if len(txs) > prt.MaxTxsPerBlock {
		return txs[:prt.MaxTxsPerBlock]
	}

	return txs
}

func (p *Mempool) DelTx(txId prt.Hash) {
	p.mu.Lock()
	defer p.mu.Unlock()

	strTxId := utils.HashToString(txId)

	delete(p.transactions, strTxId)
}

// Mempool 초기화
func (p *Mempool) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.transactions = make(map[string]*Transaction)
}
