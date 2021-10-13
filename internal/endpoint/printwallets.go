package endpoint

import (
	"encoding/json"
	"io"
	"net/http"
)

func (e *EndPoint) handlePrintWallets(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(e.wallets.GetAllSeeds(), "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}
