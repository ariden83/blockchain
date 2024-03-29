package explorer

import (
	"net/http"
	"net/url"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-session/session"
	"go.uber.org/zap"

	"github.com/ariden83/blockchain/cmd/web/internal/auth/classic"
	"github.com/ariden83/blockchain/cmd/web/internal/decoder"
	"github.com/ariden83/blockchain/cmd/web/internal/ip"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

type loginData struct {
	*FrontData
	Success    bool
	Paraphrase string
}

func (e *Explorer) loginPage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	if authorized {
		http.Redirect(rw, r, defaultPageLogged, http.StatusFound)
		return
	}

	frontData := loginData{
		FrontData: e.frontData(rw, r).
			JS([]string{
				"https://www.google.com/recaptcha/api.js?render=" + e.cfg.ReCaptcha.SiteKey,
				"/static/login/login.js?v0.0.30",
				"/static/cipher.js?v0.0.8",
			}).
			Css([]string{
				"/static/login/login.css?0.0.2",
			}).
			Title("login"),
		Success:    false,
		Paraphrase: decoder.GetPrivateKey(),
	}

	e.ExecuteTemplate(rw, r, "login", frontData)
}

func (e *Explorer) logoutPage(rw http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), rw, r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	store.Delete(sessionLabelAccessToken)
	store.Delete(sessionLabelRefreshToken)
	store.Save()
	rw.Header().Set("Location", defaultPageLogin)
	rw.WriteHeader(http.StatusFound)
}

func (e *Explorer) authorizePage(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parm := r.Form
	if parm == nil {
		parm = url.Values{}
	}

	parm.Add("grant_type", "client_credentials")
	parm.Add("client_id", e.auth.API[classic.Name].Config().ClientID)
	parm.Add("client_secret", e.auth.API[classic.Name].Config().ClientSecret)
	parm.Add("scope", "all")
	parm.Add("response_type", "token")

	r.Form = parm

	req, err := e.authServer.ValidationAuthorizeRequest(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		// err := srv.redirectError(w, req, err)}
		return
	}

	// user authorization
	address, err := e.authServer.UserAuthorizationHandler(rw, r)
	if err != nil {
		//return s.redirectError(w, req, err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		// err := srv.redirectError(w, req, err)}
		return
	} else if address == "" {
		rw.Header().Set("Location", defaultPageLogin)
		return
	}
	req.UserID = address

	// specify the scope of authorization
	if fn := e.authServer.AuthorizeScopeHandler; fn != nil {
		scope, err := fn(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		} else if scope != "" {
			req.Scope = scope
		}
	}

	// specify the expiration time of access token
	if fn := e.authServer.AccessTokenExpHandler; fn != nil {
		exp, err := fn(rw, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		req.AccessTokenExp = exp
	}

	ti, err := e.authServer.GetAuthorizeToken(ctx, req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the redirect URI is empty, the default domain provided by the client is used.
	if req.RedirectURI == "" {
		client, err := e.authServer.Manager.GetClient(ctx, req.ClientID)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		req.RedirectURI = client.GetDomain()
	}

	data := e.authServer.GetAuthorizeData(req.ResponseType, ti)

	/*  outputJSON(data) */
	store, err := session.Start(r.Context(), rw, r)
	if err != nil {
		return
	}
	store.Set(sessionLabelAccessToken, data["access_token"].(string))
	store.Set(sessionLabelRefreshToken, data["refresh_token"].(string))
	store.Save()

	rw.Header().Set("Location", defaultPageLogged)
	rw.WriteHeader(http.StatusFound)
}

type postLoginAPIBodyReq struct {
	Mnemonic  string `json:"mnemonic"`
	Password  string `json:"password"`
	Recaptcha string `json:"recaptcha"`
}

type postLoginAPIBodyRes struct {
	Status string `json:"status"`
}

const (
	sessionLabelUserID       string = "LoggedInUserID"
	sessionLabelAccessToken  string = "LoggedAccessToken"
	sessionLabelRefreshToken string = "LoggedRefreshToken"
	sessionCurrentBlock      string = "currentBlock"
)

const (
	defaultPageLogged string = "/wallet"
	defaultPageLogin  string = "/login"
)

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
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
// Responses:
//
//	default: genericError
//	    200: postLoginAPIResp
//	    401: genericError
//	    404: genericError
//	    412: genericError
//	    500: genericError
func (e *Explorer) loginAPI(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	if authorized {
		http.Redirect(rw, r, defaultPageLogged, http.StatusFound)
		return
	}

	req := &postLoginAPIBodyReq{}
	logCTX := e.logCTX("loginAPI")
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

	if req.Password == "" || req.Mnemonic == "" {
		logCTX.Error("missing fields")
		e.JSONfail(pkgErr.ErrMissingFields, rw)
		return
	}

	password, err := decoder.Decrypt(req.Password, decoder.GetPrivateKey())
	if err != nil {
		logCTX.Error("fail to decode password")
		e.JSONfail(err, rw)
		return
	}

	mnemonic, err := decoder.Decrypt(req.Mnemonic, decoder.GetPrivateKey())
	if err != nil {
		logCTX.Error("fail to decode mnemonic")
		e.JSONfail(err, rw)
		return
	}

	ip, err := ip.User(r)
	if err != nil {
		logCTX.Warn("fail to get user ip", zap.Error(err))
	}

	if e.reCaptcha != nil {
		if valid := e.reCaptcha.Verify(req.Recaptcha, ip); !valid {
			logCTX.Warn("fail to verify captcha")
			e.JSONfail(err, rw)
			return
		}
	}

	wallet, err := e.model.GetWallet(r.Context(), mnemonic, password)
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
	store.Set(sessionLabelUserID, string(wallet.PrivKey))
	store.Save()

	e.JSON(rw, postLoginAPIBodyRes{"ok"})
}

func (e *Explorer) tokenAPI(rw http.ResponseWriter, r *http.Request) {
	if err := e.refreshToken(rw, r); err != nil {
		_, statusCode, _ := e.authServer.GetErrorData(err)
		http.Error(rw, err.Error(), statusCode)
		return
	}
	rw.Header().Set("Location", defaultPageLogged)
	rw.WriteHeader(http.StatusFound)
}

func (e *Explorer) refreshToken(rw http.ResponseWriter, r *http.Request) error {
	store, err := session.Start(r.Context(), rw, r)
	if err != nil {
		return err
	}
	refreshToken, ok := store.Get(sessionLabelRefreshToken)
	if !ok {
		return err
	}

	parm := r.Form
	if parm == nil {
		parm = url.Values{}
	}
	parm.Add("refresh_token", refreshToken.(string))
	parm.Add("grant_type", oauth2.Refreshing.String())
	parm.Add("client_id", e.auth.API[classic.Name].Config().ClientID)
	parm.Add("client_secret", e.auth.API[classic.Name].Config().ClientSecret)
	parm.Add("scope", "all")

	r.Form = parm

	ctx := r.Context()

	gt, tgr, err := e.authServer.ValidationTokenRequest(r)
	if err != nil {
		return err
	}

	ti, err := e.authServer.GetAccessToken(ctx, gt, tgr)
	if err != nil {
		return err
	}

	data := e.authServer.GetTokenData(ti)
	store.Set(sessionLabelAccessToken, data["access_token"].(string))
	store.Set(sessionLabelRefreshToken, data["refresh_token"].(string))
	store.Save()
	return nil
}
