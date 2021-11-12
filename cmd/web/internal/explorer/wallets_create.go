package explorer

import (
	"encoding/json"
	"github.com/ariden83/blockchain/internal/wallet"
	"net/http"
)

type apiParamInput struct{}
type apiParamOutput struct {
	wallet.Seed
}

type walletsCreateData struct {
	PageTitle string
	Phrase    string
}

func (e *Explorer) walletsCreate(rw http.ResponseWriter, r *http.Request) {
	var (
		params apiParamInput = apiParamInput{}
		path   string        = "/wallet"
		data   apiParamOutput
	)

	body, err := e.model.Post(path, params)
	if err != nil {
		return
	}

	json.NewDecoder(body).Decode(&data)

	frontData := walletsCreateData{"Wallets", data.Mnemonic}

	templates.ExecuteTemplate(rw, "wallets_index", frontData)
}
