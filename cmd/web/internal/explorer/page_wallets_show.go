package explorer

import (
	"net/http"

	"github.com/ariden83/blockchain/cmd/web/internal/utils"
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

func walletsShowPage(rw http.ResponseWriter, r *http.Request) {
	address := utils.GetRoute(r, "address")
	outputs := []*UnspTxOutput{}
	balance := uint(0)

	data := walletsShowData{
		PageTitle:     "Show Wallet",
		Address:       address,
		Balance:       balance,
		UnspTxOutputs: outputs,
	}
	templates.ExecuteTemplate(rw, "wallets_show", data)
}
