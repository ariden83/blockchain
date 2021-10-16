package endpoint

import (
	"net/http"
)

func (e *EndPoint) handlePrintWallets(w http.ResponseWriter, r *http.Request) {
	allSeeds := e.wallets.GetAllPublicSeeds()
	e.respondWithJSON(w, r, http.StatusBadRequest, allSeeds)
}
