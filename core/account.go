package core

import (
	"fmt"

	prt "github.com/abcfe/abcfe-node/protocol"
)

type Account struct {
	Address   prt.Address `json:"address"`   // 계정 주소
	Status    string      `json:"status"`    // 계정 상태 (STAKING, VALIDATING, PROPOSERING, JAILED)
	Balance   uint64      `json:"balance"`   // 잔액
	CreatedAt int64       `json:"createdAt"` // 계정 최초 활성화 시간
	UpdatedAt int64       `json:"updateAt"`  // 마지막 수정 시간
}

// func setAccount(address Address) *Account {

// }

func NewAccountTx(address prt.Address, tx prt.Hash) {
	// address 기반으로 불러오고 PrefixAccountTx
	// 있다면 추가하고 UpdatedAt 업데이트
	// 없다면 json 형태로 생성 setAccount
	// key:value로 저장
	panic("Not Developed Yet")
}

func NewAccountTxTo(address prt.Address, tx prt.Hash) {
	// address 기반으로 불러오고 PrefixAccountTxTo
	// 있다면 추가하고 UpdatedAt 업데이트
	// 없다면 json 형태로 생성 setAccount
	// key:value로 저장
	panic("Not Developed Yet")
}

func NewAccountTxFrom(address prt.Address, tx prt.Hash) {
	// address 기반으로 불러오고 PrefixAccountTxFrom
	// 있다면 추가하고 UpdatedAt 업데이트
	// 없다면 json 형태로 생성 setAccount
	// key:value로 저장
	panic("Not Developed Yet")
}

func UpdAccountBalance(address prt.Address) {
	// address 기반으로 불러오고 PrefixAccount
	// UTXO 기반으로 balance 계산 // utxo.go에 있을거임
	// 있다면 balance와 UpdatedAt 업데이트
	// 없다면 json 형태로 생성 setAccount
	// key:value로 저장
	panic("Not Developed Yet")
}

func UpdAccountStatus(address prt.Address, status string) {
	// address 기반으로 불러오고 PrefixAccount
	// 있다면 추가하고 UpdatedAt 업데이트
	// 없다면 json 형태로 생성 setAccount
	// key:value로 저장
	panic("Not Developed Yet")
}

func (p *BlockChain) GetBalance(address prt.Address) (uint64, error) {
	utxoList, err := p.GetUtxoList(address, false)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	balance := p.CalBalanceUtxo(utxoList)
	return balance, nil
}
