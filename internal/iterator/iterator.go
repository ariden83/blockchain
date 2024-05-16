// Package iterator implements a blockchain iterator.
package iterator

import (
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/utils"
)

// BlockChainIterator represent a blockchain iterator.
type BlockChainIterator struct {
	CurrentHash []byte
	Persistence persistenceadapter.Adapter
	// Database    *badger.DB
}

// New returns a new iterator takes our BlockChain struct and returns it as a BlockChainIterator.
func New(p persistenceadapter.Adapter) *BlockChainIterator {
	iterator := BlockChainIterator{
		CurrentHash: p.LastHash(),
		Persistence: p,
	}

	return &iterator
}

// Next can get the next blockchain iterator.
func (b *BlockChainIterator) Next() (*blockchain.Block, error) {
	val, err := b.Persistence.GetCurrentHashSerialize(b.CurrentHash)
	if err != nil {
		return nil, err
	}

	block := &blockchain.Block{}
	if err := utils.Deserialize(val, block); err != nil {
		return nil, err
	}

	b.CurrentHash = block.PrevHash
	return block, nil
}
