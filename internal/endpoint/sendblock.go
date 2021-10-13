package endpoint

import (
	"encoding/json"
	"fmt"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/handle"
	"github.com/ariden83/blockchain/internal/transactions"
	"github.com/ariden83/blockchain/internal/utils"
	"github.com/davecgh/go-spew/spew"
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

func (e *EndPoint) handleSendBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m SendBlockInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	e.sendBlock(w, m)
}

func (e *EndPoint) sendBlock(w http.ResponseWriter, input SendBlockInput) {
	lastHash, index := e.getLastBlock()

	tx, err := e.transaction.New(input.From, input.To, input.Amount)
	if err == transactions.ErrNotEnoughFunds {
		io.WriteString(w, "Transaction failed, not enough funds")
		return

	} else {
		handle.Handle(err)
	}

	newBlock := blockchain.AddBlock(lastHash, index, tx)

	if blockchain.IsBlockValid(newBlock, blockchain.BlockChain[len(blockchain.BlockChain)-1]) {

		mutex.Lock()
		blockchain.BlockChain = append(blockchain.BlockChain, newBlock)
		mutex.Unlock()

		ser, err := utils.Serialize(&newBlock)
		handle.Handle(err)

		err = e.persistence.Update(newBlock.Hash, ser)
		handle.Handle(err)
		spew.Dump(blockchain.BlockChain)

	} else {
		handle.Handle(fmt.Errorf("new block is invalid"))
	}

	io.WriteString(w, "Transaction done")
}
