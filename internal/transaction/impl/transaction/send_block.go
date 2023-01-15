package transaction

import (
	"errors"
	"math/big"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/p2p/address"
	"github.com/ariden83/blockchain/internal/p2p/validator"
	"github.com/ariden83/blockchain/internal/transaction"
	pkgError "github.com/ariden83/blockchain/pkg/errors"
)

type SendBlockInput struct {
	From   []byte
	To     []byte
	Amount *big.Int
}

func (t *Transactions) SendBlock(input transaction.SendBlockInput) error {
	lastHash, index, err := t.GetLastBlock()
	if err != nil {
		return err
	}

	tx, err := t.New(input.From, input.To, input.Amount)
	if err == pkgError.ErrNotEnoughFunds {
		t.log.Info("Transaction failed, not enough funds",
			zap.Any("param", input),
			zap.String("input", "sendBlock"))
		return pkgError.ErrNotEnoughFunds
	} else if err != nil {
		return err
	}

	newBlock := blockchain.AddBlock(lastHash, index, tx)

	if !blockchain.IsBlockValid(newBlock, blockchain.BlockChain[len(blockchain.BlockChain)-1]) {
		return errors.New("new block is invalid")
	}

	mutex.Lock()
	t.event.PushBlock(validator.New(newBlock, address.IAM.CurrentAddress()))
	mutex.Unlock()

	/*mutex.Lock()
	blockchain.BlockChain = append(blockchain.BlockChain, newBlock)
	mutex.Unlock()

	ser, err := utils.Serialize(&newBlock)
	e.Handle(err)

	e.event.Push(event.BlockChain)

	err = e.persistence.Update(newBlock.Hash, ser)
	e.Handle(err)
	spew.Dump(blockchain.BlockChain)*/

	return nil
}
