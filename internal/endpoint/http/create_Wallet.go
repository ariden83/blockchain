package http

import (
	"github.com/ariden83/blockchain/internal/event"
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

	e.event.Push(event.Message{Type: event.Wallet})

	e.JSON(w, http.StatusCreated, newSeed)
}
