package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (e *EndPoint) makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/ping", e.handleGetPong).Methods(http.MethodGet)
	muxRouter.HandleFunc("/blockchain", e.handleGetBlockChain).Methods(http.MethodGet)
	muxRouter.HandleFunc("/balance", e.handleGetBalance).Methods(http.MethodPost)
	muxRouter.HandleFunc("/block", e.handleCreateBlock).Methods(http.MethodPost)
	muxRouter.HandleFunc("/send", e.handleSendBlock).Methods(http.MethodPost)
	muxRouter.HandleFunc("/wallets", e.handleGetWallets).Methods(http.MethodGet)
	muxRouter.HandleFunc("/wallet", e.handleCreateWallet).Methods(http.MethodPost)
	muxRouter.HandleFunc("/wallet", e.handleGetWallet).Methods(http.MethodGet)
	muxRouter.HandleFunc("/address", e.handleGetServersAddress).Methods(http.MethodGet)

	muxRouter.Use(defaultHeader)
	muxRouter.Use(e.MetricsMiddleware)

	return muxRouter
}
