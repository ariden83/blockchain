package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/pkg/api"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

// handleGetWallet
func (e *EndPoint) handleGetWallet(rw http.ResponseWriter, r *http.Request) {
	req := &api.GetWalletInput{}

	log := e.log.With(zap.String("input", "myWallet"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		e.log.Error("invalid params", zap.Error(err))
		return
	}

	if string(req.Mnemonic) == "" || string(req.Password) == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error(err.Error())
		return
	}

	seed, err := e.wallets.Seed(req.GetMnemonic(), req.GetPassword())
	if err != nil {
		return
	}

	e.JSON(rw, http.StatusCreated, api.GetWalletOutput{
		Address: seed.Address,
		PubKey:  seed.PubKey,
		PrivKey: seed.PrivKey,
	})
}
