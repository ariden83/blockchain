package http

import (
	"github.com/ariden83/blockchain/internal/p2p/address"
	"net/http"
)

func (e *EndPoint) handleGetServersAddress(w http.ResponseWriter, r *http.Request) {
	e.respondWithJSON(w, http.StatusCreated, address.GetCurrentAddress())
}
