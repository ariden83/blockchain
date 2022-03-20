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
	From   string   `json:"from"`
	To     string   `json:"to"`
	Amount *big.Int `json:"amount"`
}

func (e *EndPoint) handleSendBlock(rw http.ResponseWriter, r *http.Request) {
	req := &SendBlockInput{}

	log := e.log.With(zap.String("input", "sendBlock"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		e.Handle(err)
	}

	if req.From == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error("Missing from param", zap.Error(err))
		e.Handle(err)
	}

	if req.To == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error("Missing to param", zap.Error(err))
		e.Handle(err)
	}

	e.transaction.SendBlock(transactions.SendBlockInput{
		From:   []byte(req.From),
		To:     []byte(req.To),
		Amount: req.Amount,
	})
}
