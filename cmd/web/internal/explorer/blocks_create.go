package explorer

import (
	"net/http"
)

func blocksCreate(rw http.ResponseWriter, r *http.Request) {

	http.Redirect(rw, r, "/blocks", http.StatusFound)
}
