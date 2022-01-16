package explorer

import (
	"net/http"
)

func (e *Explorer) inscriptionPage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	if authorized {
		http.Redirect(rw, r, "/wallet", http.StatusFound)
		return
	}
	frontData := frontData{"Seed creation", authorized}

	templates.ExecuteTemplate(rw, "inscription", frontData)
}

type postinscriptionAPIBodyReq struct{}

type postinscriptionAPIBodyRes struct {
	Address  string `json:"address"`
	PubKey   string `json:"pubkey"`
	Mnemonic string `json:"mnemonic"`
}

// postinscriptionAPIResp
//
// swagger:response postinscriptionAPIResp
// nolint
type postinscriptionAPIResp struct {
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
		postinscriptionAPIBodyRes
	} `json:"body"`
}

// postinscriptionAPIReq Params for method POST
//
// swagger:parameters postinscriptionAPIReq
// nolint
type postinscriptionAPIReq struct {
	// the items for this order
	//
	// in: body
	postinscriptionAPIBodyReq postinscriptionAPIBodyReq
	// X-Request-Id
	// in: header
	// required: true
	XRequestID string `json:"X-Request-Id"`
	// X-Token
	// in: header
	// required: true
	XToken string `json:"X-Token"`
}

// inscriptionAPI swagger:route POST /api/registration inscriptionAPI postinscriptionAPIReq
//
// POST inscriptionAPI
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
//        200: postinscriptionAPIResp
//        401: genericError
//        404: genericError
//        412: genericError
//        500: genericError
func (e *Explorer) inscriptionAPI(rw http.ResponseWriter, r *http.Request) {

	wallet, err := e.model.CreateWallet(r.Context())
	if err != nil {
		e.fail(http.StatusNotFound, err, rw)
		return
	}

	e.JSON(rw, postinscriptionAPIBodyRes{
		Address:  wallet.Address,
		PubKey:   wallet.PubKey,
		Mnemonic: wallet.Mnemonic,
	})
}
