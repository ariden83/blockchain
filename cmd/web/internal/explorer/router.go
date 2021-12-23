package explorer

func (e *Explorer) loadRoutes() {
	e.router.HandleFunc("/", homePage).Methods("GET")
	e.router.HandleFunc("/wallet/login", e.loginPage).Methods("GET")
	e.router.HandleFunc("/wallet/create", e.walletsCreatePage).Methods("GET")
	e.router.HandleFunc("/404", notFoundPage).Methods("GET")

	e.router.HandleFunc("/blocks", blocksIndexPage).Methods("GET")
	e.router.HandleFunc("/blocks/{hash:[0-9a-f]+}", blocksShowPage).Methods("GET")
	e.router.HandleFunc("/blocks", blocksCreatePage).Methods("POST")
	e.router.HandleFunc("/blocks/mine", blocksMinePage).Methods("GET")

	e.router.HandleFunc("/transactions/{id:[0-9a-f]+}", txsShowPage).Methods("GET")

	e.router.HandleFunc("/wallets", walletsIndexPage).Methods("GET")
	e.router.HandleFunc("/wallets/server", walletsServerPage).Methods("GET")
	e.router.HandleFunc("/wallets/{address:[0-9a-f]+}", walletsShowPage).Methods("GET")
}

func (e *Explorer) loadAPIRoutes() {
	s := e.router.PathPrefix("/api").Subrouter()
	s.HandleFunc("/login", e.loginAPI).Methods("POST")
}
