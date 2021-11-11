package explorer

import (
	"net/http"
)

type walletsServerData struct {
	PageTitle     string
	Address       string
	Balance       uint
	UnspTxOutputs []*UnspTxOutput
}

func walletsServer(rw http.ResponseWriter, r *http.Request) {
	address := ""
	outputs := []*UnspTxOutput{}
	balance := uint(0)

	data := walletsServerData{
		PageTitle:     "Show Wallet",
		Address:       address,
		Balance:       balance,
		UnspTxOutputs: outputs,
	}
	templates.ExecuteTemplate(rw, "wallets_show", data)
}
