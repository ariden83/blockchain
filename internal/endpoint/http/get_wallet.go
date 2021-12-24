package http

import (
	"github.com/ariden83/blockchain/internal/wallet"
	"go.uber.org/zap"
	"net/http"
)

type GetWalletInput struct {
	Mnemonic string `json:"mnemonic"`
}

type GetWalletOutput struct {
	Address string `json:"address"`
	PubKey  string `json:"public_key"`
}

func (e *EndPoint) handleGetWallet(rw http.ResponseWriter, r *http.Request) {
	req := &GetWalletInput{}

	log := e.log.With(zap.String("input", "myWallet"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		return
	}

	keys := wallet.GetKeys(req.Mnemonic)
	e.JSONRes(rw, http.StatusCreated, GetWalletOutput{
		Address: keys.Address,
		PubKey:  keys.PubKey,
	})
}
