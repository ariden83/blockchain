package endpoint

import (
	"net/http"
)

func (e *EndPoint) handlePrintWallets(w http.ResponseWriter, r *http.Request) {
	allSeeds := e.wallets.GetAllSeeds()
	e.respondWithJSON(w, r, http.StatusBadRequest, allSeeds)
}
