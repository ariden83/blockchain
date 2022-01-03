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

func (e *Explorer) walletsCreatePage(rw http.ResponseWriter, r *http.Request) {
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

	templates.ExecuteTemplate(rw, "inscription", frontData)
}

type postRegistrationAPIBodyReq struct{}

type postRegistrationAPIBodyRes struct {
	Address  string `json:"address"`
	PubKey   string `json:"pubkey"`
	Mnemonic string `json:"mnemonic"`
}

// postRegistrationAPIResp
//
// swagger:response postRegistrationAPIResp
// nolint
type postRegistrationAPIResp struct {
	// Content-Length
	// in: header
	// required: true
	ContentLength string `json:"Content-Length"`
	// Content-Type
	// in: header
	// required: true
	ContentType string `json:"Content-Type"`
	// X-Request-Id
	// in: header
	// required: true
	XRequestID string `json:"X-Request-Id"`
	// corps of Response
	// in: body
	Body struct {
		postRegistrationAPIBodyRes
	} `json:"body"`
}

// postRegistrationAPIReq Params for method POST
//
// swagger:parameters postRegistrationAPIReq
// nolint
type postRegistrationAPIReq struct {
	// the items for this order
	//
	// in: body
	postRegistrationAPIBodyReq postRegistrationAPIBodyReq
	// X-Request-Id
	// in: header
	// required: true
	XRequestID string `json:"X-Request-Id"`
	// X-Token
	// in: header
	// required: true
	XToken string `json:"X-Token"`
}

// registrationAPI swagger:route POST /api/registration registrationAPI postRegistrationAPIReq
//
// POST registrationAPI
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
// Responses:
//    default: genericError
//        200: postRegistrationAPIResp
//        401: genericError
//        404: genericError
//        412: genericError
//        500: genericError
func (e *Explorer) registrationAPI(rw http.ResponseWriter, r *http.Request) {
	wallet, err := e.model.CreateWallet(r.Context())
	if err != nil {
		e.fail(http.StatusNotFound, err, rw)
		return
	}

	e.JSON(rw, postRegistrationAPIBodyRes{
		Address:  wallet.Address,
		PubKey:   wallet.PubKey,
		Mnemonic: wallet.Mnemonic,
	})
}
