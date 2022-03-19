package transactions

import (
	"errors"
	"math/big"

	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/p2p/address"
	"github.com/ariden83/blockchain/internal/p2p/validation"
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

	block, err := utils.DeserializeBlock(serializeBloc)
	if err != nil {
		return lastHash, nil, err
	}

	return lastHash, block.Index, nil
}

type WriteBlockInput struct {
	PubKey     string `json:"key"`
	PrivateKey string `json:"private"`
}

func (t *Transactions) WriteBlock(p WriteBlockInput) (*blockchain.Block, error) {
	lastHash, index, err := t.GetLastBlock()
	if err != nil {
		return nil, err
	}
	//mutex.Lock()
	cbtx := t.CoinBaseTx(p.PubKey, "")
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
