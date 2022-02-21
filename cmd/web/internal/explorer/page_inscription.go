package explorer

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/go-session/session"

	"github.com/ariden83/blockchain/cmd/web/internal/decoder"
	"github.com/ariden83/blockchain/cmd/web/internal/ip"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

type inscriptionData struct {
	*FrontData
	Success    bool
	Paraphrase string
}

func (e *Explorer) inscriptionPage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	if authorized {
		http.Redirect(rw, r, "/wallet", http.StatusFound)
		return
	}

	frontData := inscriptionData{
		FrontData: e.frontData(rw, r).
			JS([]string{
				"https://www.google.com/recaptcha/api.js?render=" + e.cfg.ReCaptcha.SiteKey,
				"/static/inscription/inscription.js?v0.0.16",
				"/static/qrious.min.js",
				"/static/cipher.js?v0.0.3",
			}).
			Css([]string{
				"/static/inscription/inscription.css?0.0.2",
			}).
			Title("inscription"),
		Success:    false,
		Paraphrase: decoder.GetPrivateKey(),
	}

	e.ExecuteTemplate(rw, r, "inscription", frontData)
}

type postInscriptionAPIBodyReq struct {
	Password  string `json:"password"`
	Recaptcha string `json:"recaptcha"`
}

type postInscriptionAPIBodyRes struct {
	Status string `json:"status"`
	Seed   string `json:"seed"`
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
		e.JSONfail(pkgErr.ErrAlreadyConnected, rw)
		return
	}

	req := &postInscriptionAPIBodyReq{}
	logCTX := e.logCTX("inscriptionAPI")

	if err := e.decodeBody(rw, logCTX, r.Body, req); err != nil {
		logCTX.Error("fail to decode body", zap.Error(err))
		e.JSONfail(pkgErr.ErrMissingFields, rw)
		return
	}
	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			logCTX.Error("fail to parse form", zap.Error(err))
			e.JSONfail(pkgErr.ErrMissingFields, rw)
			return
		}
	}

	if req.Password == "" {
		logCTX.Error("missing password", zap.String("password", req.Password))
		e.JSONfail(pkgErr.ErrMissingPassword, rw)
		return
	}

	ip, err := ip.User(r)
	if err != nil {
		logCTX.Warn("fail to get user ip", zap.Error(err))
	}
	if e.reCaptcha != nil {
		if valid := e.reCaptcha.Verify(req.Recaptcha, ip); !valid {
			logCTX.Warn("fail to verify captcha", zap.String("captcha", req.Recaptcha))
			e.JSONfail(err, rw)
			return
		}
	}

	password, err := decoder.Decrypt(req.Password, decoder.GetPrivateKey())
	if err != nil {
		logCTX.Error("fail to decode password", zap.String("password", req.Password))
		e.JSONfail(err, rw)
		return
	}

	wallet, err := e.model.CreateWallet(r.Context(), password)
	if err != nil {
		e.JSONfail(err, rw)
		return
	}

	store, err := session.Start(r.Context(), rw, r)
	if err != nil {
		logCTX.Error("fail to start session", zap.Error(err))
		e.JSONfail(pkgErr.ErrInternalError, rw)
		return
	}
	store.Set(sessionLabelUserID, wallet.PubKey)
	store.Save()

	mnemonic, err := decoder.Encrypt([]byte(wallet.Mnemonic), decoder.GetPrivateKey())
	if err != nil {
		logCTX.Error("fail to Encrypt mnemonic", zap.String("mnemonic", string(wallet.Mnemonic)))
		e.JSONfail(pkgErr.ErrInternalError, rw)
		return
	}

	e.JSON(rw, postInscriptionAPIBodyRes{
		Status: "ok",
		Seed:   mnemonic,
	})
}

func (e *Explorer) inscriptionValidateAPI(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	if authorized {
		http.Redirect(rw, r, defaultPageLogged, http.StatusFound)
		return
	}

	logCTX := e.logCTX("inscriptionValidateAPI")

	store, err := session.Start(r.Context(), rw, r)
	if err != nil {
		logCTX.Error("fail to start session", zap.Error(err))
		e.JSONfail(pkgErr.ErrInternalError, rw)
		return
	}

	pubKey, ok := store.Get(sessionLabelUserID)
	if !ok {
		e.JSONfail(pkgErr.ErrInternalError, rw)
		return
	}

	output, err := e.model.ValidWallet(r.Context(), []byte(pubKey.(string)))
	if err != nil {
		e.JSONfail(err, rw)
		return
	} else if !output.Valid {
		e.JSONfail(pkgErr.ErrInternalError, rw)
		return
	}

	e.JSON(rw, postInscriptionAPIBodyRes{
		Status: "ok",
	})
}
