package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/event"
	"github.com/ariden83/blockchain/pkg/api"
)

func (e *EndPoint) handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	req := &api.CreateWalletInput{}

	log := e.log.With(zap.String("input", "createBlock"))
	if err := e.decodeBody(w, log, r.Body, req); err != nil {
		return
	}

	seed, err := e.wallets.Create(req.Password)
	if err != nil {
		e.log.Error("Fail to create wallet", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	e.event.Push(event.Message{Type: event.Wallet})

	e.JSON(w, http.StatusCreated, &api.CreateWalletOutput{
		Address:  seed.Address,
		PubKey:   seed.PubKey,
		PrivKey:  seed.PrivKey,
		Mnemonic: seed.Mnemonic,
	})
}
