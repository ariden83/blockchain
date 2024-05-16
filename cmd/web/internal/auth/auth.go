package auth

import (
	"net/http"

	"golang.org/x/oauth2"

	"github.com/ariden83/blockchain/cmd/web/config"
	"github.com/ariden83/blockchain/cmd/web/internal/auth/classic"
	"github.com/ariden83/blockchain/cmd/web/internal/auth/google"
)

type IOAuth interface {
	GenerateStateOauthCookie(http.ResponseWriter) string
	GetUserData(code string) ([]byte, error)
	Config() *oauth2.Config
}

type Auth struct {
	API map[string]IOAuth
}

func New(options ...func(*Auth)) *Auth {
	e := &Auth{
		API: map[string]IOAuth{},
	}
	for _, o := range options {
		o(e)
	}
	return e
}

func WithGoogleAPI(cfg config.OAuthConfig) func(*Auth) {
	return func(e *Auth) {
		if cfg.Enable {
			e.API[google.Name] = google.New(cfg)
		}
	}
}

func WithClassic(cfg config.OAuthConfig) func(*Auth) {
	return func(e *Auth) {
		if cfg.Enable {
			e.API[classic.Name] = classic.New(cfg)
		}
	}
}
