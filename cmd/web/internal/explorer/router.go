package explorer

import "net/http"

func (e *Explorer) loadRoutes() {
	e.router.HandleFunc("/", homePage).Methods(http.MethodGet)
	e.router.HandleFunc("/wallet/login", e.loginPage).Methods(http.MethodGet)
	e.router.HandleFunc("/wallet/create", e.walletsCreatePage).Methods(http.MethodGet)
	e.router.HandleFunc("/404", notFoundPage).Methods(http.MethodGet)

	e.router.HandleFunc("/blocks", blocksIndexPage).Methods(http.MethodGet)
	e.router.HandleFunc("/blocks/{hash:[0-9a-f]+}", blocksShowPage).Methods(http.MethodGet)
	e.router.HandleFunc("/blocks", blocksCreatePage).Methods("POST")
	e.router.HandleFunc("/blocks/mine", blocksMinePage).Methods(http.MethodGet)

	e.router.HandleFunc("/transactions/{id:[0-9a-f]+}", txsShowPage).Methods(http.MethodGet)

	e.router.HandleFunc("/wallets", walletsIndexPage).Methods(http.MethodGet)
	e.router.HandleFunc("/wallets/server", walletsServerPage).Methods(http.MethodGet)
	e.router.HandleFunc("/wallets/{address:[0-9a-f]+}", walletsShowPage).Methods(http.MethodGet)
}

func (e *Explorer) loadAPIRoutes() {
	s := e.router.PathPrefix("/api").Subrouter()
	s.HandleFunc("/login", e.loginAPI).Methods(http.MethodPost)
}
