package explorer

import (
	"fmt"
	"net/http"
)

func (e *Explorer) loadFileServer() {
	fmt.Println(e.cfg.StaticDir)
	// fileServer := http.FileServer(http.Dir(e.cfg.StaticDir))

	fs := http.FileServer(http.Dir("./static"))
	e.router.Handle(e.cfg.StaticRoute, fs)
	// e.router.Handle(e.cfg.StaticRoute, http.StripPrefix(e.cfg.StaticDir, fileServer))
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
	e.router.HandleFunc("/wallets/create", e.walletsCreate).Methods("GET")
	e.router.HandleFunc("/wallets/login", e.walletsLoginForm).Methods("GET")
	e.router.HandleFunc("/wallets/server", walletsServer).Methods("GET")
	e.router.HandleFunc("/wallets/{address:[0-9a-f]+}", walletsShow).Methods("GET")

}
