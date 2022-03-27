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

// Iterator takes our BlockChain struct and returns it as a BlockCHainIterator struct
func New(p persistence.IPersistence) *BlockChainIterator {
	iterator := BlockChainIterator{
		CurrentHash: p.LastHash(),
		Persistence: p,
	}

	return &iterator
}

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
