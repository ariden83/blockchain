package http

import (
	"go.uber.org/zap"
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

func (e *EndPoint) handleGetBalance(rw http.ResponseWriter, r *http.Request) {
	req := &getBalanceInput{}

	log := e.log.With(zap.String("input", "getBalance"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		return
	}

	balance := e.transaction.FindUserBalance(req.PubKey)
	tokensSend := e.transaction.FindUserTokensSend(req.PubKey)
	tokensReceived := e.transaction.FindUserTokensReceived(req.PubKey)

	e.JSON(rw, http.StatusOK, getBalanceOutput{
		Address:       req.PubKey,
		Balance:       balance,
		TotalReceived: tokensReceived,
		TotalSent:     tokensSend,
	})
}
