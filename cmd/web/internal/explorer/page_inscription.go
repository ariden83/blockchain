package explorer

import (
	"errors"
	"github.com/ariden83/blockchain/cmd/web/internal/decoder"
	"go.uber.org/zap"
	"net/http"

	"github.com/go-session/session"

	"github.com/ariden83/blockchain/cmd/web/internal/ip"
)

type inscriptionData struct {
	*FrontData
	Success    bool
	Paraphrase string
}

const passwordKey string = "~NB8CcOL#J!H?|Yr"

func (e *Explorer) inscriptionPage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	if authorized {
		http.Redirect(rw, r, "/wallet", http.StatusFound)
		return
	}

	frontData := inscriptionData{
		FrontData: e.frontData(rw, r).
			JS([]string{
				"https://cdnjs.cloudflare.com/ajax/libs/crypto-js/4.0.0/crypto-js.min.js",
				"https://www.google.com/recaptcha/api.js?render=" + e.cfg.ReCaptcha.SiteKey,
				"/static/inscription/inscription.js?v0.0.12",
			}).
			Css([]string{
				"/static/inscription/inscription.css?0.0.0",
			}).
			Title("inscription"),
		Success:    false,
		Paraphrase: passwordKey,
	}

	e.ExecuteTemplate(rw, r, "inscription", frontData)
}

type postInscriptionAPIBodyReq struct {
	Cipher    string `json:"cipher"`
	IV        string `json:"iv"`
	Recaptcha string `json:"recaptcha"`
}

type postInscriptionAPIBodyRes struct {
	Status string `json:"status"`
}

// postInscriptionAPIResp
//
// swagger:response postInscriptionAPIResp
// nolint
type postInscriptionAPIResp struct {
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
		postInscriptionAPIBodyRes
	} `json:"body"`
}

// postInscriptionAPIReq Params for method POST
//
// swagger:parameters postInscriptionAPIReq
// nolint
type postInscriptionAPIReq struct {
	// the items for this order
	//
	// in: body
	postInscriptionAPIBodyReq postInscriptionAPIBodyReq
	// X-Request-Id
	// in: header
	// required: true
	XRequestID string `json:"X-Request-Id"`
	// X-Token
	// in: header
	// required: true
	XToken string `json:"X-Token"`
}

// inscriptionAPI swagger:route POST /api/registration inscriptionAPI postInscriptionAPIReq
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
//        200: postInscriptionAPIResp
//        401: genericError
//        404: genericError
//        412: genericError
//        500: genericError
func (e *Explorer) inscriptionAPI(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	if authorized {
		http.Redirect(rw, r, defaultPageLogged, http.StatusFound)
		return
	}

	req := &postInscriptionAPIBodyReq{}
	log := e.log.With(zap.String("input", "oauthClassicInscription"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		e.fail(http.StatusPreconditionFailed, err, rw)
		return
	}
	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			e.fail(http.StatusInternalServerError, err, rw)
			return
		}
	}

	if req.Cipher == "" || req.IV == "" {
		e.fail(http.StatusPreconditionFailed, errors.New("missing fields"), rw)
		return
	}

	ip, err := ip.User(r)
	if e.reCaptcha != nil {
		if valid := e.reCaptcha.Verify(req.Recaptcha, ip); !valid {
			http.Error(rw, "fail to verify capcha", http.StatusPreconditionFailed)
			return
		}
	}

	password, err := decoder.Password(req.Cipher, req.IV, passwordKey)
	if err != nil {
		http.Error(rw, "fail to decode password", http.StatusPreconditionFailed)
		return
	}

	wallet, err := e.model.CreateWallet(r.Context(), password)
	if err != nil {
		e.fail(http.StatusNotFound, err, rw)
		return
	}

	store, err := session.Start(r.Context(), rw, r)
	if err != nil {
		e.fail(http.StatusNotFound, err, rw)
		return
	}
	store.Set(sessionLabelUserID, wallet.Address)
	store.Save()

	e.JSON(rw, postInscriptionAPIBodyRes{"ok"})
}
