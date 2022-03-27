package http

import (
	"net/http"
)

type getPongOutput struct {
	Ping string `json:"ping"`
}

func (e *EndPoint) handleGetPong(rw http.ResponseWriter, _ *http.Request) {
	e.JSON(rw, http.StatusOK, getPongOutput{
		Ping: "pong",
	})
}
