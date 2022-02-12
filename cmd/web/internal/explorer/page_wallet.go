package explorer

import (
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
	"github.com/go-session/session"
	"go.uber.org/zap"
	"net/http"
)

type UnspTxOutput struct {
	TxId   string `json:"transactionId"`
	Index  uint   `json:"index"`
	Amount uint   `json:"amount"`
}

type walletsData struct {
	*FrontData
	Address       string
	Balance       string
	TotalReceived string
	TotalSent     string
	UnSpTxOutputs []*UnspTxOutput
}

func (e *Explorer) walletPage(rw http.ResponseWriter, r *http.Request) {
	logCTX := e.logCTX("walletPage")
	store, err := session.Start(r.Context(), rw, r)
	if err != nil {
		logCTX.Error("fail to start session", zap.Error(err))
		e.fail(pkgErr.ErrInternalError, rw)
		return
	}
	accessToken, _ := store.Get(sessionLabelAccessToken)
	token, err := e.authServer.Manager.LoadAccessToken(r.Context(), accessToken.(string))
	if err != nil {
		logCTX.Error("fail to load access token", zap.Error(err))
		e.fail(pkgErr.ErrInternalError, rw)
		return
	}

	wallet, err := e.model.GetBalance(r.Context(), token.GetUserID(), token.GetUserID())
	if err != nil {
		logCTX.Error("fail to get balance", zap.Error(err))
		e.fail(pkgErr.ErrInternalError, rw)
		return
	}

	data := walletsData{
		Address:       wallet.Address,
		Balance:       wallet.Balance,
		TotalReceived: wallet.TotalReceived,
		TotalSent:     wallet.TotalSent,
		FrontData: &FrontData{
			PageTitle:    e.metadata.Title + " - " + "Wallet page",
			Authentified: true,
			Menus:        getMenus(),
			Javascripts:  []string{},
		},
	}

	templates.ExecuteTemplate(rw, "wallet", data)
}
