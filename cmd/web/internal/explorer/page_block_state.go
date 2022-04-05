package explorer

import (
	"net/http"

	"github.com/gorilla/websocket"

	pkgErr "github.com/ariden83/blockchain/pkg/errors"
)

type State struct {
	State   int
	Message string
}

// @doc https://gist.github.com/owulveryck/57d8c2469fd1f8a840747b064c50ff4e
// @doc https://github.com/gorilla/websocket/blob/master/examples/echo/server.go
func (e *Explorer) stateBlockAPI(rw http.ResponseWriter, r *http.Request) {
	store, userID, err := e.getUserID(rw, r)
	if err != nil {
		e.fail(pkgErr.ErrInternalError, rw)
		return
	}

	if e.ws == nil {
		e.fail(pkgErr.ErrInternalError, rw)
		return
	}
	c, err := e.ws.Upgrade(rw, r, nil)
	if err != nil {
		e.fail(pkgErr.ErrInternalError, rw)
		return
	}
	defer c.Close()

	_, ok := store.Get(sessionCurrentBlock)
	if ok {
		err = c.WriteMessage(websocket.TextMessage, []byte("minage"))
		if err != nil {
			e.fail(pkgErr.ErrInternalError, rw)
			return
		}

	} else {
		_, err = e.model.CreateBlock(r.Context(), userID)
		if err != nil {
			e.JSONfail(err, rw)
			return
		}
		store.Set(sessionCurrentBlock, "minage")

		err = c.WriteMessage(websocket.TextMessage, []byte("minage"))
		if err != nil {
			e.fail(pkgErr.ErrInternalError, rw)
			return
		}
	}

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			e.fail(pkgErr.ErrInternalError, rw)
			break
		}
		if mt != websocket.TextMessage {
			e.fail(pkgErr.ErrNotImplemented, rw)
			break
		}

		err = c.WriteMessage(mt, msg)
		if err != nil {
			e.fail(pkgErr.ErrInternalError, rw)
			break
		}
	}
}
