package iterator

import (
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/handle"
	"github.com/ariden83/blockchain/internal/persistence"
	"github.com/ariden83/blockchain/internal/utils"
)

type BlockChainIterator struct {
	CurrentHash []byte
	Persistence *persistence.Persistence
	// Database    *badger.DB
}

func (b *BlockChainIterator) Next() *blockchain.Block {
	val, err := b.Persistence.GetCurrentHashSerialize(b.CurrentHash)
	handle.Handle(err)
	block, err := utils.Deserialize(val)
	handle.Handle(err)
	b.CurrentHash = block.PrevHash
	return block
}
