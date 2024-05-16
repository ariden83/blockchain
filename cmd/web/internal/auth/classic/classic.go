package classic

import (
	"net/http"

	"golang.org/x/oauth2"

	"github.com/ariden83/blockchain/cmd/web/config"
)

const Name string = "classic"

type ClassicOAuth struct {
	config *oauth2.Config
}

func New(cfg config.OAuthConfig) *ClassicOAuth {
	return &ClassicOAuth{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Scopes:       cfg.Scopes,
			Endpoint: oauth2.Endpoint{
				TokenURL: "/api/token",
				AuthURL:  "/api/authorize",
			},
			RedirectURL: "/api/oauth2",
		},
	}
}

func (c *ClassicOAuth) Config() *oauth2.Config {
	return c.config
}

func (ClassicOAuth) GenerateStateOauthCookie(http.ResponseWriter) string {
	return ""
}

func (c *ClassicOAuth) GetUserData(string) ([]byte, error) {
	return []byte{}, nil
}
