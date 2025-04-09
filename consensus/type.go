package consensus

import (
	prt "github.com/abcfe/abcfe-node/protocol"
)

// 검증자 인터페이스
type ValidatorInfo interface {
	GetAddress() prt.Address
	GetPubKey() []byte
	GetVotingPower() uint64
	GetActiveStat() bool
}
