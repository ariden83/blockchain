package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/pkg/api"
)

func (e *EndPoint) handleGetBalance(rw http.ResponseWriter, r *http.Request) {
	req := &api.GetBalanceInput{}

	log := e.log.With(zap.String("input", "getBalance"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		return
	}

	balance := e.transaction.FindUserBalance(req.PrivKey)
	tokensSend := e.transaction.FindUserTokensSend(req.PrivKey)
	tokensReceived := e.transaction.FindUserTokensReceived(req.PrivKey)

	e.JSON(rw, http.StatusOK, api.GetBalanceOutput{
		Address:       e.wallets.GetUserAddress(req.PrivKey),
		Balance:       balance.String(),
		TotalReceived: tokensReceived.String(),
		TotalSent:     tokensSend.String(),
	})
}
