package consensus

import (
	prt "github.com/abcfe/abcfe-node/protocol"
)

type Validator struct {
	Address     prt.Address `json:"address"`
	PublicKey   []byte      `json:"publicKey"`
	VotingPower uint64      `json:"votingPower"`
	IsActive    bool        `json:"isActive"`
}

// Block Consensus Data 직렬화

// Block Consesnus Data 역직렬화

func (v *Validator) GetAddress() prt.Address {
	return v.Address
}

func (v *Validator) GetPubKey() []byte {
	return v.PublicKey
}

func (v *Validator) GetVotingPower() uint64 {
	return v.VotingPower
}

func (v *Validator) GetActiveStat() bool {
	return v.IsActive
}

func (v *Validator) SignBlock() {

}

func (v *Validator) ValidateBlock() {

}
