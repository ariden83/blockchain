package explorer

import (
	"net/http"
)

func blocksCreatePage(rw http.ResponseWriter, r *http.Request) {

	http.Redirect(rw, r, "/blocks", http.StatusFound)
}
