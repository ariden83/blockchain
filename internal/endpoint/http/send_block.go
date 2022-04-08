package http

import (
	"math/big"
	"net/http"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/internal/transactions"
	"github.com/ariden83/blockchain/pkg/api"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

func (e *EndPoint) handleSendBlock(rw http.ResponseWriter, r *http.Request) {
	req := &api.SendBlockInput{}

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
		From:   []byte(req.GetFrom()),
		To:     []byte(req.GetTo()),
		Amount: new(big.Int),
	})

	e.JSON(rw, http.StatusCreated, api.SendBlockOutput{})
}
