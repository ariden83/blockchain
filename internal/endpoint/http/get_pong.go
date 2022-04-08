package http

import (
	"net/http"

	protoAPI "github.com/ariden83/blockchain/pkg/api"
)

func (e *EndPoint) handleGetPong(rw http.ResponseWriter, _ *http.Request) {
	e.JSON(rw, http.StatusOK, protoAPI.Pong{
		Message: "pong",
	})
}
