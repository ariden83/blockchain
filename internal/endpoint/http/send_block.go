package http

import (
	"fmt"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/p2p/address"
	"github.com/ariden83/blockchain/internal/p2p/validation"
	"github.com/ariden83/blockchain/internal/transactions"
	"go.uber.org/zap"
	"io"
	"math/big"
	"net/http"
)

// Message takes incoming JSON payload for writing heart rate
type SendBlockInput struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Amount *big.Int `json:"amount"`
}

func (e *EndPoint) handleSendBlock(rw http.ResponseWriter, r *http.Request) {
	req := &SendBlockInput{}

	log := e.log.With(zap.String("input", "sendBlock"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		return
	}

	e.sendBlock(rw, *req)
}

func (e *EndPoint) sendBlock(w http.ResponseWriter, input SendBlockInput) {
	lastHash, index := e.getLastBlock()

	tx, err := e.transaction.New(input.From, input.To, input.Amount)
	if err == transactions.ErrNotEnoughFunds {
		e.log.Info("Transaction failed, not enough funds",
			zap.Any("param", input),
			zap.String("input", "sendBlock"))

		if _, err := io.WriteString(w, "Transaction failed, not enough funds"); err != nil {
			e.log.Error("fail to write string", zap.Error(err))
		}
		return

	} else {
		e.Handle(err)
	}

	newBlock := blockchain.AddBlock(lastHash, index, tx)

	if blockchain.IsBlockValid(newBlock, blockchain.BlockChain[len(blockchain.BlockChain)-1]) {

		mutex.Lock()
		e.event.PushBlock(validation.New(newBlock, address.GetCurrentAddress()))
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

	} else {
		e.Handle(fmt.Errorf("new block is invalid"))
	}

	if _, err := io.WriteString(w, "Transaction done"); err != nil {
		e.log.Error("fail to write string", zap.Error(err))
	}
}
