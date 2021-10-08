package endpoint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetBalanceInput struct {
	Address string `json:"address"`
}

func (e *EndPoint) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	var m GetBalanceInput

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	balance := 0
	UTXOs := e.transaction.FindUTXO(m.Address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	io.WriteString(w, fmt.Sprintf("Balance of %s: %d\n", m.Address, balance))
}
