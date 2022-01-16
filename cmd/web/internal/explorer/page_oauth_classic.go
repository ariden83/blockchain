package explorer

import (
	"errors"
	"github.com/ariden83/blockchain/cmd/web/internal/auth/classic"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-session/session"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

func (e *Explorer) loginPage(rw http.ResponseWriter, r *http.Request) {
	_, authorized := e.authorized(rw, r)
	if authorized {
		http.Redirect(rw, r, "/wallet", http.StatusFound)
		return
	}
	data := frontData{"Wallets connexion", authorized}
	templates.ExecuteTemplate(rw, "login", data)
}

type postLoginAPIBodyReq struct {
	Mnemonic string `json:"mnemonic"`
}

type postLoginAPIBodyRes struct {
	Address string `json:"address"`
	PubKey  string `json:"pubkey"`
}

const (
	sessionLabelUserID       string = "LoggedInUserID"
	sessionLabelAccessToken  string = "LoggedAccessToken"
	sessionLabelRefreshToken string = "LoggedRefreshToken"
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
	_, authorized := e.authorized(rw, r)
	if authorized {
		http.Redirect(rw, r, "/wallet", http.StatusFound)
		return
	}
	req := &postLoginAPIBodyReq{}
	log := e.log.With(zap.String("input", "oauthClassicLogin"))
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

	wallet, err := e.model.GetWallet(r.Context(), req.Mnemonic)
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

	rw.Header().Set("Location", "/api/authorize?grant_type=client_credentials"+
		"&client_id="+e.auth.API[classic.Name].Config().ClientID+
		"&client_secret="+e.auth.API[classic.Name].Config().ClientSecret+
		"&scope=all"+
		"&response_type=token")

	rw.WriteHeader(http.StatusFound)
}

func (e *Explorer) authorizeAPI(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := e.authServer.ValidationAuthorizeRequest(r)
	if err != nil {
		e.fail(http.StatusUnauthorized, err, rw)
		return
	}

	// user authorization
	userID, err := e.authServer.UserAuthorizationHandler(rw, r)
	if err != nil {
		e.fail(http.StatusUnauthorized, err, rw)
		return
	} else if userID == "" {
		e.fail(http.StatusUnauthorized, err, rw)
		return
	}
	req.UserID = userID

	// specify the scope of authorization
	if fn := e.authServer.AuthorizeScopeHandler; fn != nil {
		scope, err := fn(rw, r)
		if err != nil {
			e.fail(http.StatusUnauthorized, err, rw)
			return
		} else if scope != "" {
			req.Scope = scope
		}
	}

	// specify the expiration time of access token
	if fn := e.authServer.AccessTokenExpHandler; fn != nil {
		exp, err := fn(rw, r)
		if err != nil {
			e.fail(http.StatusUnauthorized, err, rw)
			return
		}
		req.AccessTokenExp = exp
	}

	ti, err := e.authServer.GetAuthorizeToken(ctx, req)
	if err != nil {
		e.fail(http.StatusUnauthorized, err, rw)
		return
	}

	// If the redirect URI is empty, the default domain provided by the client is used.
	if req.RedirectURI == "" {
		client, err := e.authServer.Manager.GetClient(ctx, req.ClientID)
		if err != nil {
			e.fail(http.StatusUnauthorized, err, rw)
			return
		}
		req.RedirectURI = client.GetDomain()
	}

	data := e.authServer.GetAuthorizeData(req.ResponseType, ti)

	/*  outputJSON(data) */
	store, err := session.Start(r.Context(), rw, r)
	if err != nil {
		e.fail(http.StatusUnauthorized, err, rw)
		return
	}
	store.Set(sessionLabelAccessToken, data["access_token"].(string))
	store.Set(sessionLabelRefreshToken, data["refresh_token"].(string))
	store.Save()

	rw.Header().Set("Location", "/protected")
	rw.WriteHeader(http.StatusFound)
}

func (e *Explorer) tokenAPI(rw http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), rw, r)
	if err != nil {
		e.fail(http.StatusUnauthorized, err, rw)
		return
	}
	refreshToken, ok := store.Get(sessionLabelRefreshToken)
	if !ok {
		e.fail(http.StatusUnauthorized, err, rw)
		return
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
		_, statusCode, _ := e.authServer.GetErrorData(err)
		e.fail(statusCode, errors.New("fail to valid token request"), rw)
		return
	}

	ti, err := e.authServer.GetAccessToken(ctx, gt, tgr)
	if err != nil {
		_, statusCode, _ := e.authServer.GetErrorData(err)
		e.fail(statusCode, errors.New("fail to valid token request"), rw)
		return
	}

	data := e.authServer.GetTokenData(ti)
	store.Set(sessionLabelAccessToken, data["access_token"].(string))
	store.Set(sessionLabelRefreshToken, data["refresh_token"].(string))
	store.Save()

	rw.Header().Set("Location", "/protected")
	rw.WriteHeader(http.StatusFound)
}
