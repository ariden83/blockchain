package explorer

import (
	"github.com/ariden83/blockchain/cmd/web/internal/utils"
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

func walletPage(rw http.ResponseWriter, r *http.Request) {

	address := utils.GetRoute(r, "address")
	outputs := []*UnspTxOutput{}
	balance := uint(0)

	data := walletsShowData{
		PageTitle:     "your wallet",
		Address:       address,
		Balance:       balance,
		UnspTxOutputs: outputs,
	}

	templates.ExecuteTemplate(rw, "wallet", data)
}
