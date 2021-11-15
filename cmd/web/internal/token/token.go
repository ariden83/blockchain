package token

import (
	"github.com/ariden83/blockchain/cmd/web/internal/config"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

type Token struct {
	cfg config.Token
}

func New(cfg config.Token) *Token {
	return &Token{
		cfg: cfg,
	}
}

func (t *Token) CreateToken(publicKeyID string) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["public_key_id"] = publicKeyID
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(t.cfg.SecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}
