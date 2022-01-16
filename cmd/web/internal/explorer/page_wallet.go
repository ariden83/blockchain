package explorer

import (
	"github.com/go-session/session"
	"net/http"
)

type UnspTxOutput struct {
	TxId   string `json:"transactionId"`
	Index  uint   `json:"index"`
	Amount uint   `json:"amount"`
}

type walletsShowData struct {
	PageTitle     string
	Address       string
	Balance       uint
	UnspTxOutputs []*UnspTxOutput
}

func (e *Explorer) walletPage(rw http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), rw, r)
	if err != nil {
		e.fail(http.StatusInternalServerError, err, rw)
		return
	}
	accessToken, _ := store.Get(sessionLabelAccessToken)
	token, err := e.authServer.Manager.LoadAccessToken(r.Context(), accessToken.(string))
	if err != nil {
		e.fail(http.StatusInternalServerError, err, rw)
		return
	}

	outputs := []*UnspTxOutput{}
	balance := uint(0)

	wallet, err := e.model.GetWallet(r.Context(), token.GetUserID())
	if err != nil {
		e.fail(http.StatusNotFound, err, rw)
		return
	}

	data := walletsShowData{
		PageTitle:     "your wallet",
		Address:       wallet.Address,
		Balance:       balance,
		UnspTxOutputs: outputs,
	}

	templates.ExecuteTemplate(rw, "wallet", data)
}
