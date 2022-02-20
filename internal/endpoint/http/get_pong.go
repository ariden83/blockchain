package http

import (
	"math/big"
	"net/http"
)

type getPongOutput struct {
	Address            string
	Balance            *big.Int
	TotalReceived      *big.Int
	TotalSent          *big.Int
	UnconfirmedBalance *big.Int
}

func (e *EndPoint) handleGetPOng(rw http.ResponseWriter, _ *http.Request) {
	e.JSON(rw, http.StatusOK, getPongOutput{})
}
