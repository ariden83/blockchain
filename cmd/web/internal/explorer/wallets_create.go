package explorer

import (
	"encoding/json"
	"fmt"
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
	Token     string
}

func (e *Explorer) walletsCreate(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("***************************************** walletsCreate")
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

	token, err := e.token.CreateToken(data.PubKey)
	if err != nil {
		templates.ExecuteTemplate(rw, "error", Error{http.StatusUnauthorized, err})
		return
	}

	frontData := walletsCreateData{"Wallets", data.Mnemonic, token}

	templates.ExecuteTemplate(rw, "wallets_create", frontData)
}
