package http

import (
	"net/http"

	"github.com/ariden83/blockchain/internal/p2p/address"
)

func (e *EndPoint) handleGetServersAddress(w http.ResponseWriter, r *http.Request) {
	e.JSON(w, http.StatusCreated, address.IAM.CurrentAddress())
}
