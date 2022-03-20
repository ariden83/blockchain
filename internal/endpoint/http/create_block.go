package http

import (
	"net/http"

	"go.uber.org/zap"

	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

// Message takes incoming JSON payload for writing heart rate
type CreateBlockInput struct {
	PrivateKey string `json:"private"`
}

// handleCreateBlock takes JSON payload as an input for heart rate (BPM)
func (e *EndPoint) handleCreateBlock(rw http.ResponseWriter, r *http.Request) {
	req := &CreateBlockInput{}

	log := e.log.With(zap.String("input", "createBlock"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		e.Handle(err)
	}

	if req.PrivateKey == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error("Empty private key", zap.Error(err))
		e.Handle(err)
	}

	newBlock, err := e.transaction.WriteBlock([]byte(req.PrivateKey))
	e.Handle(err)

	e.JSON(rw, http.StatusCreated, newBlock)
}
