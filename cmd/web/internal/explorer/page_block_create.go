package explorer

import (
	"net/http"

	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

func (e *Explorer) createBlockAPI(rw http.ResponseWriter, r *http.Request) {
	store, userID, err := e.getUserID(rw, r)
	if err != nil {
		e.fail(pkgErr.ErrInternalError, rw)
		return
	}

	_, ok := store.Get(sessionCurrentBlock)
	if ok {
		e.fail(pkgErr.ErrInternalError, rw)
		return
	}

	if _, err = e.model.CreateBlock(r.Context(), userID); err != nil {
		e.JSONfail(err, rw)
		return
	}

	e.JSON(rw, postLoginAPIBodyRes{"ok"})
}
