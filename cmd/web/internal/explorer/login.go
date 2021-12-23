package explorer

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
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
	PubKey  string `json:"publickey"`
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
	var req postLoginAPIBodyReq

	r.Body = http.MaxBytesReader(rw, r.Body, 1048)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		e.log.Error("fail to decode input", zap.Error(err))
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		e.log.Error(msg, zap.Error(err))
		http.Error(rw, msg, http.StatusBadRequest)
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
