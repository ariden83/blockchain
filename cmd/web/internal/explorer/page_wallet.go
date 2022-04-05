package explorer

import (
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
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
	_, userID, err := e.getUserID(rw, r)
	if err != nil {
		e.fail(pkgErr.ErrInternalError, rw)
		return
	}

	wallet, err := e.model.GetBalance(r.Context(), userID)
	if err != nil {
		logCTX.Error("fail to get balance", zap.Error(err))
		e.fail(pkgErr.ErrInternalError, rw)
		return
	}

	frontData := walletsData{
		FrontData: e.frontData(rw, r).
			JS([]string{
				"/static/wallet/wallet.js?v0.0.46",
			}).
			Css([]string{
				"/static/wallet/wallet.css?0.0.6",
			}).
			Title("wallet page"),
		Address:       wallet.Address,
		Balance:       wallet.Balance,
		TotalReceived: wallet.TotalReceived,
		TotalSent:     wallet.TotalSent,
	}

	e.ExecuteTemplate(rw, r, "wallet", frontData)
}
