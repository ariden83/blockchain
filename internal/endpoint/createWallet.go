package endpoint

import (
	"go.uber.org/zap"
	"net/http"
)

func (e *EndPoint) handleCreateWallet(w http.ResponseWriter, r *http.Request) {

	newSeed, err := e.wallets.Create()
	if err != nil {
		e.log.Error("Fail to create wallet", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	e.respondWithJSON(w, r, http.StatusCreated, newSeed)
}
