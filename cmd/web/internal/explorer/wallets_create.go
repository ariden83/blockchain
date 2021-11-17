package explorer

import (
	//	"encoding/json"
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
	/*	var (
			params    apiParamInput = apiParamInput{}
			path      string        = "/wallet"
			data      apiParamOutput
			pageTitle string = "Seed creation"
		)

		body, err := e.model.Post(path, params)
		if err != nil {
			templates.ExecuteTemplate(rw, "error", Error{http.StatusUnauthorized, err, pageTitle})
			return
		}

		json.NewDecoder(body).Decode(&data)

		token, err := e.token.CreateToken(data.PubKey)
		if err != nil {
			templates.ExecuteTemplate(rw, "error", Error{http.StatusUnauthorized, err, pageTitle})
			return
		}
		frontData := walletsCreateData{pageTitle, data.Mnemonic, token}
	*/
	frontData := walletsCreateData{"Seed creation", "eihf iefhiehfi eifh iehf eifhiehfih ehifhiehf eifhiehf", "ozijefojzeiofhioef"}

	templates.ExecuteTemplate(rw, "wallets_create", frontData)
}
