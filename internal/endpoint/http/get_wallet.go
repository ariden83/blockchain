package http

import (
	"encoding/json"
	"github.com/ariden83/blockchain/internal/wallet"
	"go.uber.org/zap"
	"io"
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
	var (
		req GetWalletInput
		err error
	)

	log := e.log.With(zap.String("input", "myWallet"))

	r.Body = http.MaxBytesReader(rw, r.Body, 1048)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err = dec.Decode(&req); err != nil {
		log.Error("fail to decode input", zap.Error(err))
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err = dec.Decode(&struct{}{}); err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		log.Error(msg, zap.Error(err))
		http.Error(rw, msg, http.StatusBadRequest)
		return
	}

	keys := wallet.GetKeys(req.Mnemonic)
	e.respondWithJSON(rw, http.StatusCreated, GetWalletOutput{
		Address: keys.Address,
		PubKey:  keys.PubKey,
	})
}
