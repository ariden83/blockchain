package validator

import (
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/p2p/address"
)

type Validator struct {
	Block    blockchain.Block
	Accepted map[string]bool
	Refused  map[string]bool
	Servers  []string
}

// New creates a new validator.
func New(block blockchain.Block, servers []string) Validator {
	return Validator{
		Block:    block,
		Servers:  servers,
		Accepted: map[string]bool{},
		Refused:  map[string]bool{},
	}
}

func (v *Validator) IsAcceptedByMajority() bool {
	if len(v.Accepted) > (len(v.Servers) / 2) {
		return true
	}
	return false
}

func (v *Validator) IsRefusedByMajority() bool {
	if len(v.Refused) > (len(v.Servers) / 2) {
		return true
	}
	return false
}

func (v *Validator) Accept() {
	v.Accepted[address.IAM.Address()] = true
}

func (v *Validator) Refuse() {
	v.Refused[address.IAM.Address()] = true
}
