package transactions

import (
	"fmt"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/p2p/address"
	"github.com/ariden83/blockchain/internal/p2p/validation"
	"github.com/ariden83/blockchain/internal/transactions/trace"
	pkgError "github.com/ariden83/blockchain/pkg/errors"
)

func (t *Transactions) WriteBlock(privateKey []byte) (*blockchain.Block, error) {
	lastHash, index, err := t.GetLastBlock()
	if err != nil {
		return nil, err
	}

	t.trace.Push(blockchain.NextID(index).String(), trace.Minage)
	fmt.Println(fmt.Sprintf("*************************************** WriteBlock 2"))

	//mutex.Lock()
	cbtx := t.CoinBaseTx(privateKey)
	fmt.Println(fmt.Sprintf("*************************************** WriteBlock 3"))
	cbtx.SetID()

	newBlock := blockchain.AddBlock(lastHash, index, cbtx)
	fmt.Println(fmt.Sprintf("*************************************** WriteBlock 4"))
	t.trace.Push(newBlock.Index.String(), trace.Create)
	fmt.Println(fmt.Sprintf("*************************************** WriteBlock 5"))
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
		return nil, pkgError.ErrCreatedBlockIsInvalid
	}

	return &newBlock, nil
}
