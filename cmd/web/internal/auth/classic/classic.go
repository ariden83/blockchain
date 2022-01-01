package classic

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/ariden83/blockchain/cmd/web/config"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const Name string = "classic"

type ClassicOAuth struct {
	config *oauth2.Config
	urlAPI string
}

func New(cfg config.OAuthConfig) *ClassicOAuth {
	return &ClassicOAuth{
		config: &oauth2.Config{
			ClientID:     os.Getenv("CLASSIC_OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("CLASSIC_OAUTH_CLIENT_SECRET"),
			Scopes:       cfg.Scopes,
			Endpoint: oauth2.Endpoint{
				TokenURL: "/token",
				AuthURL:  "/authorize",
			},
			RedirectURL: "/auth/callback",
		},
		urlAPI: cfg.URLAPI,
	}
}

func (c *ClassicOAuth) Config() *oauth2.Config {
	return c.config
}

func (ClassicOAuth) GenerateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func (c *ClassicOAuth) GetUserData(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.
	token, err := c.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(c.urlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
