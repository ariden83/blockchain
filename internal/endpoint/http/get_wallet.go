package http

import (
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

type GetWalletInput struct {
	Mnemonic string `json:"mnemonic"`
	Password string `json:"password"`
}

type GetWalletOutput struct {
	Address  string `json:"address"`
	PubKey   string `json:"public_key"`
	Mnemonic string `json:"mnemonic"`
}

func (e *EndPoint) handleGetWallet(rw http.ResponseWriter, r *http.Request) {
	req := &GetWalletInput{}

	log := e.log.With(zap.String("input", "myWallet"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		e.log.Error("invalid params", zap.Error(err))
		return
	}

	if req.Mnemonic == "" || req.Password == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error(err.Error(), zap.String("mnemonic", req.Mnemonic))
		return
	}

	seed, err := e.wallets.GetSeed([]byte(req.Mnemonic), []byte(req.Password))
	if err != nil {
		return
	}

	e.JSON(rw, http.StatusCreated, GetWalletOutput{
		Address:  seed.Address,
		PubKey:   seed.PubKey,
		Mnemonic: seed.Mnemonic,
	})
}
