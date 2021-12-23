package http

import (
	"github.com/ariden83/blockchain/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func (e *EndPoint) makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/blockchain", e.handleGetBlockChain).Methods("GET")
	muxRouter.HandleFunc("/balance", e.handleGetBalance).Methods("POST")
	muxRouter.HandleFunc("/block", e.handleCreateBlock).Methods("POST")
	muxRouter.HandleFunc("/send", e.handleSendBlock).Methods("POST")
	muxRouter.HandleFunc("/wallets", e.handleGetWallets).Methods("GET")
	muxRouter.HandleFunc("/wallet", e.handleCreateWallet).Methods("POST")
	muxRouter.HandleFunc("/wallet", e.handleGetWallet).Methods("GET")
	muxRouter.HandleFunc("/address", e.handleGetServersAddress).Methods("GET")

	muxRouter.Use(middleware.DefaultHeader)
	muxRouter.Use(e.MetricsMiddleware)

	return muxRouter
}
