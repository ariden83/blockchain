package endpoint

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"math/big"
	"net/http"
)

type getBalanceInput struct {
	Address string `json:"address"`
	PubKey  string `json:"key"`
}

type getBalanceOutput struct {
	Address            string
	Balance            *big.Int
	TotalReceived      *big.Int
	TotalSent          *big.Int
	UnconfirmedBalance *big.Int
}

func (e *EndPoint) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	var input getBalanceInput

	r.Body = http.MaxBytesReader(w, r.Body, 1048)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&input)
	if err != nil {
		e.log.Error("Request body must only contain a single JSON object", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		e.log.Error("Request body must only contain a single JSON object")
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	balance := e.transaction.FindUserBalance(input.PubKey)
	tokensSend := e.transaction.FindUserTokensSend(input.PubKey)
	tokensReceived := e.transaction.FindUserTokensReceived(input.PubKey)

	e.respondWithJSON(w, r, http.StatusOK, getBalanceOutput{
		Address:       input.PubKey,
		Balance:       balance,
		TotalReceived: tokensReceived,
		TotalSent:     tokensSend,
	})
}
