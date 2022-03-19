package transactions

import (
	"errors"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/p2p/address"
	"github.com/ariden83/blockchain/internal/p2p/validation"
)

type WriteBlockInput struct {
	PubKey     []byte
	PrivateKey []byte
}

func (t *Transactions) WriteBlock(p WriteBlockInput) (*blockchain.Block, error) {
	lastHash, index, err := t.GetLastBlock()
	if err != nil {
		return nil, err
	}
	//mutex.Lock()
	cbtx := t.CoinBaseTx(p.PubKey, p.PrivateKey)
	cbtx.SetID()

	newBlock := blockchain.AddBlock(lastHash, index, cbtx)
	//mutex.Unlock()

	if blockchain.IsBlockValid(newBlock, blockchain.BlockChain[len(blockchain.BlockChain)-1]) {

		mutex.Lock()
		t.event.PushBlock(validation.New(newBlock, address.GetCurrentAddress()))
		mutex.Unlock()

		/*mutex.Lock()
		blockchain.BlockChain = append(blockchain.BlockChain, newBlock)
		mutex.Unlock()

		ser, err := utils.Serialize(&newBlock)
		e.Handle(err)

		e.event.Push(event.BlockChain, "")

		err = e.persistence.Update(newBlock.Hash, ser)
		e.Handle(err)
		spew.Dump(blockchain.BlockChain)*/
	} else {
		return nil, errors.New("new block created is invalid")
	}

	return &newBlock, nil
}
