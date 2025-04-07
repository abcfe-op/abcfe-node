package consensus

import "github.com/abcfe/abcfe-node/core"

// 검증자 인터페이스
type ValidatorInfo interface {
	GetAddress() core.Address
	GetPubKey() []byte
	GetVotingPower() uint64
	GetActiveStat() bool
}
