package http

import (
	"math/big"
	"net/http"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/transactions"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

// Message takes incoming JSON payload for writing heart rate
type SendBlockInput struct {
	PubKey  string   `json:"pubKey"`
	PrivKey string   `json:"privkey"`
	To      string   `json:"to"`
	Amount  *big.Int `json:"amount"`
}

func (e *EndPoint) handleSendBlock(rw http.ResponseWriter, r *http.Request) {
	req := &SendBlockInput{}

	log := e.log.With(zap.String("input", "sendBlock"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		e.Handle(err)
	}

	if req.PubKey == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error("Empty pub key", zap.Error(err))
		e.Handle(err)
	}

	if req.To == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error("Fail to create block", zap.Error(err))
		e.Handle(err)
	}

	e.transaction.SendBlock(transactions.SendBlockInput{
		PubKey:  []byte(req.PubKey),
		PrivKey: []byte(req.PrivKey),
		To:      []byte(req.To),
		Amount:  req.Amount,
	})
}
