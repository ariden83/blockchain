package explorer

import (
	"go.uber.org/zap"
	"net/http"
)

type loginForm struct {
	PageTitle string
}

func (e *Explorer) loginPage(rw http.ResponseWriter, r *http.Request) {
	frontData := loginForm{"Wallets connexion"}
	templates.ExecuteTemplate(rw, "wallets_login_form", frontData)
}

type postLoginAPIBodyReq struct {
	Mnemonic string `json:"mnemonic"`
}

type postLoginAPIBodyRes struct {
	Address string `json:"address"`
	PubKey  string `json:"pubkey"`
}

// postLoginResp
//
// swagger:response postLoginResp
// nolint
type postLoginResp struct {
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
		postLoginAPIBodyRes
	} `json:"body"`
}

// postLoginAPIReq Params for method POST
//
// swagger:parameters postLoginAPIReq
// nolint
type postLoginAPIReq struct {
	// the items for this order
	//
	// in: body
	postLoginAPIBodyReq postLoginAPIBodyReq
	// X-Request-Id
	// in: header
	// required: true
	XRequestID string `json:"X-Request-Id"`
	// X-Token
	// in: header
	// required: true
	XToken string `json:"X-Token"`
}

// loginAPI swagger:route POST /api/login loginAPI postLoginAPIReq
//
// POST loginAPI
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
//        200: postLoginAPIResp
//        401: genericError
//        404: genericError
//        412: genericError
//        500: genericError
func (e *Explorer) loginAPI(rw http.ResponseWriter, r *http.Request) {
	req := &postLoginAPIBodyReq{}

	log := e.log.With(zap.String("input", "loginAPI"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		return
	}

	wallet, err := e.model.GetWallet(r.Context(), req.Mnemonic)
	if err != nil {
		e.fail(http.StatusNotFound, err, rw)
		return
	}

	e.resp(rw, postLoginAPIBodyRes{
		Address: wallet.Address,
		PubKey:  wallet.PubKey,
	})
}
