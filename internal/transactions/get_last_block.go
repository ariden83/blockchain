package transactions

import (
	"errors"
	"math/big"

	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/utils"
)

func (t *Transactions) GetLastBlock() ([]byte, *big.Int, error) {
	lastHash, err := t.persistence.GetLastHash()
	if err != nil {
		return lastHash, nil, err
	}

	if lastHash == nil {
		return lastHash, nil, errors.New("no hash found")
	}

	serializeBloc, err := t.persistence.GetCurrentHashSerialize(lastHash)
	if err != nil {
		return lastHash, nil, err
	}

	block := &blockchain.Block{}
	if err := utils.Deserialize(serializeBloc, block); err != nil {
		return lastHash, nil, err
	}

	return lastHash, block.Index, nil
}
