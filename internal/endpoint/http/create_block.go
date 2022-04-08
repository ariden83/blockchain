package http

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/ariden83/blockchain/pkg/api"
	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

// handleCreateBlock takes JSON payload as an input for heart rate (BPM)
func (e *EndPoint) handleCreateBlock(rw http.ResponseWriter, r *http.Request) {
	req := &api.CreateBlockInput{}

	log := e.log.With(zap.String("input", "createBlock"))
	if err := e.decodeBody(rw, log, r.Body, req); err != nil {
		e.Handle(err)
	}

	if string(req.PrivKey) == "" {
		err := pkgErr.ErrMissingFields
		e.log.Error("Empty private key", zap.Error(err))
		e.Handle(err)
	}

	go func() {
		_, err := e.transaction.WriteBlock(req.PrivKey)
		if err != nil {
			e.log.Error("invalid block", zap.Error(err))
		}
	}()

	e.JSON(rw, http.StatusProcessing, &api.CreateBlockOutput{})
}
