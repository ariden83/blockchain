package explorer

import (
	"net/http"
)

const (
	staticDir   string = "explorer/static"
	staticRoute string = "/static/"
)

func (e *Explorer) loadFileServer() {
	fileServer := http.FileServer(http.Dir(staticDir))
	e.router.Handle(staticRoute, http.StripPrefix(staticRoute, fileServer))
}

func (e *Explorer) loadRoutes() {
	e.router.HandleFunc("/", home).Methods("GET")
	e.router.HandleFunc("/404", notFound).Methods("GET")

	e.router.HandleFunc("/blocks", blocksIndex).Methods("GET")
	e.router.HandleFunc("/blocks/{hash:[0-9a-f]+}", blocksShow).Methods("GET")
	e.router.HandleFunc("/blocks", blocksCreate).Methods("POST")
	e.router.HandleFunc("/blocks/mine", blocksMine).Methods("GET")

	e.router.HandleFunc("/transactions/{id:[0-9a-f]+}", txsShow).Methods("GET")

	e.router.HandleFunc("/wallets", walletsIndex).Methods("GET")
	e.router.HandleFunc("/wallets/server", walletsServer).Methods("GET")
	e.router.HandleFunc("/wallets/{address:[0-9a-f]+}", walletsShow).Methods("GET")

	e.router.HandleFunc("/wallets/create", e.walletsCreate).Methods("GET")
}
