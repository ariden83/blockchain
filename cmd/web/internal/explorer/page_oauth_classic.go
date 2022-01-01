package explorer

import (
	"go.uber.org/zap"
	"net/http"
	"net/url"

	"github.com/go-session/session"
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
func (e *Explorer) oauthLogin(rw http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, rw, r)
	if err != nil {
		e.fail(http.StatusInternalServerError, err, rw)
		return
	}

	req := &postLoginAPIBodyReq{}

	log := e.log.With(zap.String("input", "oauthClassicLogin"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		e.fail(http.StatusPreconditionFailed, err, rw)
		return
	}

	wallet, err := e.model.GetWallet(r.Context(), req.Mnemonic)
	if err != nil {
		e.fail(http.StatusNotFound, err, rw)
		return
	}

	store.Set("LoggedInUserID", wallet.PubKey)
	store.Save()

	rw.WriteHeader(http.StatusFound)

	e.JSON(rw, postLoginAPIBodyRes{
		Address: wallet.Address,
		PubKey:  wallet.PubKey,
	})
}

func (e *Explorer) oauthCallback(rw http.ResponseWriter, r *http.Request) {}

type authBodyRes struct {
	Location string
}

func (e *Explorer) oauthHandler(rw http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, rw, r)
	if err != nil {
		e.fail(http.StatusInternalServerError, err, rw)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		rw.WriteHeader(http.StatusFound)
		e.JSON(rw, authBodyRes{
			Location: "/login",
		})
		return
	}

	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
	}

	u := new(url.URL)
	u.Path = "/authorize"
	u.RawQuery = form.Encode()
	rw.Header().Set("Location", u.String())
	rw.WriteHeader(http.StatusFound)
	store.Delete("Form")

	if v, ok := store.Get("LoggedInUserID"); ok {
		store.Set("UserID", v)
	}
	store.Save()

	e.JSON(rw, authBodyRes{
		Location: u.String(),
	})
}

func (e *Explorer) userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("UserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)
	store.Delete("UserID")
	store.Save()
	return
}
