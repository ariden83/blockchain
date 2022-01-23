package iterator

import (
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/utils"
)

type BlockChainIterator struct {
	CurrentHash []byte
	Persistence persistence.IPersistence
	// Database    *badger.DB
}

func (b *BlockChainIterator) Next() (*blockchain.Block, error) {
	val, err := b.Persistence.GetCurrentHashSerialize(b.CurrentHash)
	if err != nil {
		return nil, err
	}

	block, err := utils.DeserializeBlock(val)
	if err != nil {
		return nil, err
	}

	b.CurrentHash = block.PrevHash
	return block, nil
}

// Iterator takes our BlockChain struct and returns it as a BlockCHainIterator struct
func New(p persistence.IPersistence) *BlockChainIterator {
	iterator := BlockChainIterator{
		CurrentHash: p.LastHash(),
		Persistence: p,
	}

	return &iterator
}
