package endpoint

import (
	"encoding/json"
	"github.com/wemeetagain/go-hdwallet"
	"io"
	"net/http"
)

type getWalletInput struct {
	Mnemonic string `json:"mnemonic"`
}

type getWalletOutput struct {
	Address    string `json:"address"`
	PubKey     string `json:"publickey"`
	PrivateKey string `json:"privkey"`
}

func (e *EndPoint) handleMyWallet(w http.ResponseWriter, r *http.Request) {
	var p getWalletInput

	r.Body = http.MaxBytesReader(w, r.Body, 1048)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Get private key from mnemonic
	masterprv := hdwallet.MasterKey([]byte(p.Mnemonic))

	// Convert a private key to public key
	masterpub := masterprv.Pub()

	// Get your address
	address := masterpub.Address()
	respondWithJSON(w, r, http.StatusCreated, getWalletOutput{
		Address:    address,
		PubKey:     masterpub.String(),
		PrivateKey: masterprv.String(),
	})
}
