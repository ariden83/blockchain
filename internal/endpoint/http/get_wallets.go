package http

import (
	"net/http"
)

func (e *EndPoint) handleGetWallets(w http.ResponseWriter, r *http.Request) {
	allSeeds := e.wallets.GetSeeds()
	e.JSON(w, http.StatusBadRequest, allSeeds)
}
