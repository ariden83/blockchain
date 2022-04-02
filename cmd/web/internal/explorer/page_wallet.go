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

	accessToken, ok := store.Get(sessionLabelAccessToken)
	if !ok {
		logCTX.Error("fail to get token")
		e.fail(pkgErr.ErrInternalError, rw)
		return
	}

	token, err := e.authServer.Manager.LoadAccessToken(r.Context(), accessToken.(string))
	if err != nil {
		logCTX.Error("fail to load access token", zap.Error(err))
		e.fail(pkgErr.ErrInternalError, rw)
		return
	}

	wallet, err := e.model.GetBalance(r.Context(), token.GetUserID())
	if err != nil {
		logCTX.Error("fail to get balance", zap.Error(err))
		e.fail(pkgErr.ErrInternalError, rw)
		return
	}

	frontData := walletsData{
		FrontData: e.frontData(rw, r).
			JS([]string{
				"/static/wallet/vue-simple-progress.min.js?v0.0.2",
				"/static/wallet/wallet.js?v0.0.1",
			}).
			Css([]string{
				"/static/wallet/wallet.css?0.0.1",
			}).
			Title("wallet page"),
		Address:       wallet.Address,
		Balance:       wallet.Balance,
		TotalReceived: wallet.TotalReceived,
		TotalSent:     wallet.TotalSent,
	}

	e.ExecuteTemplate(rw, r, "wallet", frontData)
}
