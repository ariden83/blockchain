package http

import (
	"github.com/ariden83/blockchain/internal/event"
	"go.uber.org/zap"
	"net/http"
)

type CreateWalletInput struct {
	Password string `json:"password"`
}

type CreateWalletOutput struct {
	Address  string `json:"address"`
	PubKey   string `json:"public_key"`
	Mnemonic []byte `json:"mnemonic"`
}

func (e *EndPoint) handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	req := &CreateWalletInput{}

	log := e.log.With(zap.String("input", "createBlock"))
	if err := e.decodeBody(w, log, r.Body, req); err != nil {
		return
	}

	seed, err := e.wallets.Create([]byte(req.Password))
	if err != nil {
		e.log.Error("Fail to create wallet", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	e.event.Push(event.Message{Type: event.Wallet})

	e.JSON(w, http.StatusCreated, GetWalletOutput{
		Address:  seed.Address,
		PubKey:   seed.PubKey,
		Mnemonic: seed.Mnemonic,
	})
}
