package http

import (
	"net/http"

	"go.uber.org/zap"
)

func (e *EndPoint) handleGetWallets(w http.ResponseWriter, r *http.Request) {
	allSeeds, err := e.wallets.GetSeeds()
	if err != nil {
		e.log.Error("Fail to get wallets", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	e.JSON(w, http.StatusOK, allSeeds)
}
