package endpoint

import (
	"encoding/json"
	"github.com/ariden83/blockchain/internal/blockchain"
	"github.com/ariden83/blockchain/internal/handle"
	"github.com/ariden83/blockchain/internal/transactions"
	"io"
	"net/http"
)

// Message takes incoming JSON payload for writing heart rate
type SendBlockInput struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
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

	blockchain.AddBlock(lastHash, index, tx)

	io.WriteString(w, "Transaction done")
}
