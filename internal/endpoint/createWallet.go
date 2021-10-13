package endpoint

import "net/http"

func (e *EndPoint) handleCreateWallet(w http.ResponseWriter, r *http.Request) {

	newSeed := e.wallets.Create()

	respondWithJSON(w, r, http.StatusCreated, newSeed)
}
