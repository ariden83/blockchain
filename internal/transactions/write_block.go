package transactions

import (
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/event/trace"
	"github.com/ariden83/blockchain/internal/p2p/address"
	"github.com/ariden83/blockchain/internal/p2p/validation"
	pkgError "github.com/ariden83/blockchain/pkg/errors"
)

func (t *Transactions) WriteBlock(privateKey []byte) (*blockchain.Block, error) {
	lastHash, index, err := t.GetLastBlock()
	if err != nil {
		return nil, err
	}

	t.event.PushTrace(blockchain.NextID(index).String(), trace.Minage)

	//mutex.Lock()
	cbtx := t.CoinBaseTx(privateKey)
	cbtx.SetID()

	newBlock := blockchain.AddBlock(lastHash, index, cbtx)
	//mutex.Unlock()

	if blockchain.IsBlockValid(newBlock, blockchain.BlockChain[len(blockchain.BlockChain)-1]) {

		mutex.Lock()
		t.event.PushTrace(newBlock.Index.String(), trace.Create)
		t.event.PushBlock(validation.New(newBlock, address.GetCurrentAddress()))
		mutex.Unlock()
	} else {
		return nil, pkgError.ErrCreatedBlockIsInvalid
	}

	return &newBlock, nil
}
